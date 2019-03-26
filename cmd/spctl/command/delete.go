//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/oracle/speedle/cmd/spctl/client"

	"github.com/spf13/cobra"
)

var (
	deleteExample = `
		# Delete service "foo"
		spctl delete service foo
		
		# Delete all policies in service "foo"
		spctl delete policy --all --service-name=foo
		
		# Delete policy "p01" in service "foo"
		spctl delete policy p01 --service-name=foo
		
		# Delete policy "p01" in  admin policy
		spctl delete policy p01 --admin-policy
		
		# Delete function "foo"
		spctl delete function foo
		
		# Delete all functions
		spctl delete function --all`
)

func NewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete (service | policy | rolepolicy | function) (--all | NAME | ID) [--service-name=NAME]",
		Short:   "Delete one or many services | policies | role-policies",
		Example: deleteExample,
		Run:     deleteCommandFunc,
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "Delete all elements")
	cmd.Flags().StringVar(&serviceName, "service-name", "", "Service name")
	return cmd
}

func deleteCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Help()
		return
	}

	hc, err := httpClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	cli := &client.Client{
		PMSEndpoint: globalFlags.PMSEndpoint,
		HTTPClient:  hc,
	}

	switch strings.ToLower(args[0]) {
	case "service":
		if all {
			err = cli.Delete([]string{"service"}, "")
		} else {
			if len(args[1:]) == 0 {
				cmd.Help()
				return
			}
			for _, name := range args[1:] {
				err = cli.Delete([]string{"service", name}, "")
				if err != nil {
					break
				}
			}
		}
	case "policy", "rolepolicy":
		if serviceName == "" {
			cmd.Help()
			return
		}
		var kind string
		if "policy" == strings.ToLower(args[0]) {
			kind = "policy"
		} else {
			kind = "role-policy"
		}
		if all {
			err = cli.Delete([]string{"service", serviceName, kind}, "")
		} else {
			if len(args[1:]) == 0 {
				cmd.Help()
				return
			}
			for _, name := range args[1:] {
				err = cli.Delete([]string{"service", serviceName, kind, name}, "")
				if err != nil {
					break
				}
			}
		}
	case "function":
		if all {
			err = cli.Delete([]string{"function"}, "")
		} else {
			if len(args[1:]) == 0 {
				cmd.Help()
				return
			}
			for _, name := range args[1:] {
				err = cli.Delete([]string{"function", name}, "")
				if err != nil {
					break
				}
			}
		}
	default:
		cmd.Help()
		return
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("%s deleted.\n", strings.Join(args, " "))
	}
}
