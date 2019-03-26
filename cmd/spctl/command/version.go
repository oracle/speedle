//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Short:   "Print version information",
		Example: "spctl version",
		Run:     versionCommand,
	}

	return cmd
}

func versionCommand(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fmt.Println("Usage: spctl version")
		os.Exit(1)
	}
}
