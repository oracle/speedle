//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package main

import (
	"fmt"
	"os"

	"github.com/oracle/speedle/cmd/spctl/command"
)

var gitCommit string
var productVersion string
var goVersion string

func printVersionInfo() {
	if len(os.Args) != 2 || os.Args[1] != "version" {
		return
	}
	fmt.Printf("spxctl:\n")
	fmt.Printf(" Version:       %s\n", productVersion)
	fmt.Printf(" Go Version:    %s\n", goVersion)
	fmt.Printf(" Git commit:    %s\n", gitCommit)
	os.Exit(0)
}

func main() {
	printVersionInfo()
	command.Execute()
}
