package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/DuKanghub/dcert/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	challenge := "dns"

	var obtainCmd = &cli.Command{
		Name:    "obtain",
		Aliases: []string{"o"},
		Usage:   "obtain new certificate",
		Action: func(c *cli.Context) error {
			domainStr := c.Args().Get(0)
			domains := strings.Split(domainStr, ",")

			cmd.GetSSLCerts(challenge, domains)
			return nil
		},
	}

	var renewCmd = &cli.Command{
		Name:    "renew",
		Aliases: []string{"r"},
		Usage:   "renew certificate",
		Action: func(c *cli.Context) error {
			domainStr := c.Args().Get(0)
			domains := strings.Split(domainStr, ",")
			cmd.RenewCert(challenge, domains)
			return nil
		},
	}
	commands := []*cli.Command{obtainCmd, renewCmd}

	app := &cli.App{
		Name:                 "CertManager",
		Usage:                "A command line tool to obtain certificate",
		Version:              "0.0.1",
		Description:          "A command line tool to obtain certificate",
		Commands:             commands,
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
