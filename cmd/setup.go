package cmd

import (
	"os"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

const filePerm os.FileMode = 0o600

func createNonExistingFolder(path string) (err error) {
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0o700)
	}

	return
}

func setup(accountsStorage *AccountsStorage) (account *AcmeAccount, client *lego.Client) {
	privateKey := accountsStorage.GetPrivateKey(certcrypto.RSA2048)

	if accountsStorage.ExistsAccountFilePath() {
		account = accountsStorage.LoadAccount(privateKey)
	} else {
		account = &AcmeAccount{Email: accountsStorage.GetUserID(), privateKey: privateKey}
	}

	client, err := newClient(account)
	if err != nil {
		panic(err)
	}

	return
}

func newClient(user registration.User) (client *lego.Client, err error) {
	legoConfig := lego.NewConfig(user)

	// you can also add custom parameters to config
	// legoConfig.CADirURL = config.CONFIG.String("CADirURL")

	client, err = lego.NewClient(legoConfig)

	return
}
