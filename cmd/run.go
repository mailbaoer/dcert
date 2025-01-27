package cmd

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DuKanghub/dcert/config"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

const rootPathWarningMessage = `!!!! HEADS UP !!!!

Your account credentials have been saved in your Let's Encrypt
configuration directory at "%s".

You should make a secure backup of this folder now. This
configuration directory will also contain certificates and
private keys obtained from Let's Encrypt so making regular
backups of this folder is ideal.
`

func GetSSLCerts(challenge string, domains []string) {
	email := config.CONFIG.String("Email")

	accountsStorage := NewAccountsStorage(email)
	account, client := setup(accountsStorage)

	var err error
	if account.Registration == nil {
		reg, err := register(client)
		if err != nil {
			log.Fatalf("Could not complete registration\n\t%v", err)
		}
		account.Registration = reg
		if err = accountsStorage.Save(account); err != nil {
			log.Fatal(err)
		}

		fmt.Printf(rootPathWarningMessage, accountsStorage.GetRootPath())
	}

	certsStorage := NewCertificatesStorage()
	certsStorage.CreateRootFolder()

	switch challenge {
	case "dns":
		dnsProvider, err := dns01.NewDNSProviderManual()
		if err != nil {
			log.Fatal(err)
		}
		err = client.Challenge.SetDNS01Provider(
			dnsProvider,
			dns01.AddRecursiveNameservers(config.CONFIG.Strings("DNSNameservers")),
			dns01.AddDNSTimeout(60*time.Second),
		)
	case "http":
		err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "18888"))
	}

	if err != nil {
		log.Fatal(err)
	}

	cert, err := obtainCertificate(domains, client)
	if err != nil {
		log.Fatalf("Could not obtain certificates:\n\t%v", err)
	}
	certsStorage.SaveResource(cert)

}

func RenewCert(challenge string, domains []string) (err error) {
	var errStr string
	domain := domains[0]
	accountsStorage := NewAccountsStorage(config.CONFIG.String("Email"))
	account, client := setup(accountsStorage)
	if account.Registration == nil {
		log.Fatalf("Account %s is not registered. Use 'run' to register a new account.\n", account.Email)
	}
	challenge = strings.ToUpper(challenge)
	if challenge == "DNS" {
		dnsProvider, err := dns01.NewDNSProviderManual()
		if err != nil {
			log.Fatal(err)
		}
		err = client.Challenge.SetDNS01Provider(
			dnsProvider,
			dns01.AddRecursiveNameservers(config.CONFIG.Strings("DNSNameservers")),
			dns01.AddDNSTimeout(60*time.Second),
		)
	} else if challenge == "HTTP" {
		// HTTP验证：起一个http server监听18888，需保证外网能访问到。
		err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "18888"))
	}

	if err != nil {
		return err
	}
	certsStorage := NewCertificatesStorage()
	certificates, err := certsStorage.ReadCertificate(domain, ".crt")
	if err != nil {
		errStr = fmt.Sprintf("Error while loading the certificate for domain %s\n\t%v", domain, err)
		return errors.New(errStr)
	}
	cert := certificates[0]
	if !needRenewal(cert, domain, 30) {
		return nil
	}
	// This is just meant to be informal for the user.
	timeLeft := cert.NotAfter.Sub(time.Now().UTC())
	fmt.Printf("[%s] acme: Trying renewal with %d hours remaining\n", domain, int(timeLeft.Hours()))
	certDomains := certcrypto.ExtractDomains(cert)
	// 复用旧的private key
	var privateKey crypto.PrivateKey
	keyBytes, errR := certsStorage.ReadFile(domain, ".key")
	if errR != nil {
		errStr = fmt.Sprintf("Error while loading the private key for domain %s\n\t%v", domain, errR)
		return errors.New(errStr)
	}
	privateKey, errR = certcrypto.ParsePEMPrivateKey(keyBytes)
	if errR != nil {
		return errR
	}
	request := certificate.ObtainRequest{
		Domains:    merge(certDomains, domains),
		Bundle:     false,
		PrivateKey: privateKey,
	}
	certRes, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	certsStorage.SaveResource(certRes)
	return err
}

func needRenewal(x509Cert *x509.Certificate, domain string, days int) bool {
	if x509Cert.IsCA {
		log.Fatalf("[%s] Certificate bundle starts with a CA certificate", domain)
	}

	if days >= 0 {
		notAfter := int(time.Until(x509Cert.NotAfter).Hours() / 24.0)
		if notAfter > days {
			log.Printf("[%s] The certificate expires in %d days, the number of days defined to perform the renewal is %d: no renewal.",
				domain, notAfter, days)
			return false
		}
	}

	return true
}

func register(client *lego.Client) (*registration.Resource, error) {
	return client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}

func obtainCertificate(domains []string, client *lego.Client) (*certificate.Resource, error) {

	if len(domains) > 0 {
		// obtain a certificate, generating a new private key
		request := certificate.ObtainRequest{
			Domains: domains,
			Bundle:  false,
		}
		return client.Certificate.Obtain(request)
	}
	return nil, errors.New("请传入1个域名")
}

func merge(prevDomains, nextDomains []string) []string {
	for _, next := range nextDomains {
		var found bool
		for _, prev := range prevDomains {
			if prev == next {
				found = true
				break
			}
		}
		if !found {
			prevDomains = append(prevDomains, next)
		}
	}
	return prevDomains
}
