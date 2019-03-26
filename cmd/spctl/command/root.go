//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	cliName                   = "spctl"
	cliDescription            = "A command line interface for speedle"
	defaultTimeout            = 5 * time.Second
	DefaultPolicyMgmtEndPoint = "http://127.0.0.1:6733/policy-mgmt/v1/"
	DefaultAuthzCheckEndPoint = "http://127.0.0.1:6734/authz-check/v1/"
)

var (
	globalFlags = GlobalFlags{}
)

var (
	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"spctl"},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&globalFlags.PMSEndpoint, "pms-endpoint", DefaultPolicyMgmtEndPoint, "speedle policy managemnet service endpoint")
	rootCmd.PersistentFlags().DurationVar(&globalFlags.Timeout, "timeout", 5000000000, "timeout for running command")
	rootCmd.PersistentFlags().StringVar(&globalFlags.CertFile, "cert", "", "identify secure client using this TLS certificate file")
	rootCmd.PersistentFlags().StringVar(&globalFlags.KeyFile, "key", "", "identify secure client using this TLS key file")
	rootCmd.PersistentFlags().StringVar(&globalFlags.CAFile, "cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.InsecureSkipVerify, "skipverify", false, "control whether a client verifies the server's certificate chain and host name or not")

	args, _ := readConfigFile()
	for name, val := range args {
		rootCmd.PersistentFlags().Set(name, val)
	}

	rootCmd.AddCommand(
		NewGetCommand(),
		NewDeleteCommand(),
		NewCreateCommand(),
		NewConfigCommand(),
		NewDiscoverCommand(),
		NewVersionCommand(),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
