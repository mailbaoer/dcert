package cmd

import (
	"crypto"

	"github.com/go-acme/lego/v4/registration"
)

// You'll need a user or account type that implements acme.User

type AcmeAccount struct {
	Email        string
	privateKey   crypto.PrivateKey
	Registration *registration.Resource
}

func (u *AcmeAccount) GetEmail() string {
	return u.Email
}
func (u AcmeAccount) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeAccount) GetPrivateKey() crypto.PrivateKey {
	return u.privateKey
}
