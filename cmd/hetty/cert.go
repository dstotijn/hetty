package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/smallstep/truststore"
)

var certUsage = `
Usage:
    hetty cert <subcommand> [flags]

Certificate management tools.

Options:
    --help, -h  Output this usage text.

Subcommands:
    - install    Installs a certificate to the system trust store, and
                 (optionally) to the Firefox and Java trust stores.
    - uninstall  Uninstalls a certificate from the system trust store, and
                 (optionally) from the Firefox and Java trust stores.

Run ` + "`hetty cert <subcommand> --help`" + ` for subcommand specific usage instructions.

Visit https://hetty.xyz to learn more about Hetty.
`

var certInstallUsage = `
Usage:
    hetty cert install [flags]
	
Installs a certificate to the system trust store, and (optionally) to the Firefox
and Java trust stores.

Options:
    --cert         Path to certificate. (Default: "~/.hetty/hetty_cert.pem")
    --firefox      Install certificate to Firefox trust store. (Default: false)
    --java         Install certificate to Java trust store. (Default: false)
    --skip-system  Skip installing certificate to system trust store (Default: false)
    --help, -h     Output this usage text.

Visit https://hetty.xyz to learn more about Hetty.
`

var certUninstallUsage = `
Usage:
    hetty cert uninstall [flags]
	
Uninstalls a certificate from the system trust store, and (optionally) from the Firefox
and Java trust stores.

Options:
    --cert         Path to certificate. (Default: "~/.hetty/hetty_cert.pem")
    --firefox      Uninstall certificate from Firefox trust store. (Default: false)
    --java         Uninstall certificate from Java trust store. (Default: false)
    --skip-system  Skip uninstalling certificate from system trust store (Default: false)
    --help, -h     Output this usage text.

Visit https://hetty.xyz to learn more about Hetty.
`

type CertInstallCommand struct {
	config     *Config
	cert       string
	firefox    bool
	java       bool
	skipSystem bool
}

type CertUninstallCommand struct {
	config     *Config
	cert       string
	firefox    bool
	java       bool
	skipSystem bool
}

func NewCertCommand(rootConfig *Config) *ffcli.Command {
	return &ffcli.Command{
		Name: "cert",
		Subcommands: []*ffcli.Command{
			NewCertInstallCommand(rootConfig),
			NewCertUninstallCommand(rootConfig),
		},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
		UsageFunc: func(*ffcli.Command) string {
			return certUsage
		},
	}
}

func NewCertInstallCommand(rootConfig *Config) *ffcli.Command {
	cmd := CertInstallCommand{
		config: rootConfig,
	}
	fs := flag.NewFlagSet("hetty cert install", flag.ExitOnError)

	fs.StringVar(&cmd.cert, "cert", "~/.hetty/hetty_cert.pem", "Path to certificate.")
	fs.BoolVar(&cmd.firefox, "firefox", false, "Install certificate to Firefox trust store. (Default: false)")
	fs.BoolVar(&cmd.java, "java", false, "Install certificate to Java trust store. (Default: false)")
	fs.BoolVar(&cmd.skipSystem, "skip-system", false, "Skip installing certificate to system trust store (Default: false)")

	cmd.config.RegisterFlags(fs)

	return &ffcli.Command{
		Name:    "install",
		FlagSet: fs,
		Exec:    cmd.Exec,
		UsageFunc: func(*ffcli.Command) string {
			return certInstallUsage
		},
	}
}

func (cmd *CertInstallCommand) Exec(_ context.Context, _ []string) error {
	caCertFile, err := homedir.Expand(cmd.cert)
	if err != nil {
		return fmt.Errorf("failed to parse certificate filepath: %w", err)
	}

	opts := []truststore.Option{}

	if cmd.skipSystem {
		opts = append(opts, truststore.WithNoSystem())
	}

	if cmd.firefox {
		opts = append(opts, truststore.WithFirefox())
	}

	if cmd.java {
		opts = append(opts, truststore.WithJava())
	}

	if !cmd.skipSystem {
		cmd.config.logger.Info(
			"To install the certificate in the system trust store, you might be prompted for your password.")
	}

	if err := truststore.InstallFile(caCertFile, opts...); err != nil {
		return fmt.Errorf("failed to install certificate: %w", err)
	}

	cmd.config.logger.Info("Finished installing certificate.")

	return nil
}

func NewCertUninstallCommand(rootConfig *Config) *ffcli.Command {
	cmd := CertUninstallCommand{
		config: rootConfig,
	}
	fs := flag.NewFlagSet("hetty cert uninstall", flag.ExitOnError)

	fs.StringVar(&cmd.cert, "cert", "~/.hetty/hetty_cert.pem", "Path to certificate.")
	fs.BoolVar(&cmd.firefox, "firefox", false, "Uninstall certificate from Firefox trust store. (Default: false)")
	fs.BoolVar(&cmd.java, "java", false, "Uninstall certificate from Java trust store. (Default: false)")
	fs.BoolVar(&cmd.skipSystem, "skip-system", false, "Skip uninstalling certificate from system trust store (Default: false)")

	cmd.config.RegisterFlags(fs)

	return &ffcli.Command{
		Name:    "uninstall",
		FlagSet: fs,
		Exec:    cmd.Exec,
		UsageFunc: func(*ffcli.Command) string {
			return certUninstallUsage
		},
	}
}

func (cmd *CertUninstallCommand) Exec(_ context.Context, _ []string) error {
	caCertFile, err := homedir.Expand(cmd.cert)
	if err != nil {
		return fmt.Errorf("failed to parse certificate filepath: %w", err)
	}

	opts := []truststore.Option{}

	if cmd.skipSystem {
		opts = append(opts, truststore.WithNoSystem())
	}

	if cmd.firefox {
		opts = append(opts, truststore.WithFirefox())
	}

	if cmd.java {
		opts = append(opts, truststore.WithJava())
	}

	if !cmd.skipSystem {
		cmd.config.logger.Info(
			"To uninstall the certificate from the system trust store, you might be prompted for your password.")
	}

	if err := truststore.UninstallFile(caCertFile, opts...); err != nil {
		return fmt.Errorf("failed to uninstall certificate: %w", err)
	}

	cmd.config.logger.Info("Finished uninstalling certificate.")

	return nil
}
