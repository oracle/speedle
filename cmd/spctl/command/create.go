//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/oracle/speedle/api/pms"
	"github.com/oracle/speedle/cmd/spctl/client"
	"github.com/oracle/speedle/cmd/spctl/pdl"
)

var (
	pdlFileName        string
	jsonFileName       string
	command            string
	serviceType        string
	funcURL            string
	funcResultCachable bool
	funcResultTTL      int64
)

var (
	createExample = `
	    # Create an empty service with name "service1" and default service type
		spctl create service service1

		# Create an empty service with name "service1" and type "k8s"
		spctl create service service1 --service-type=k8s

		# Create a service with policies using a service definition file in json format		
		spctl create service --json-file service.json

		# Create a service with policies and role policies using a file with policies in PDL format
		spctl create service service1 --service-type=k8s --pdl-file pdl.txt 
		sample pdl-file:
		--------------------------------------------------------
		role-policies:
		grant user User1 Role1 on res1
		grant group Group1 Role2 on res2
		policies:
		grant group Administrators GET,POST,DELETE expr:/service/* if request_time > '2017-09-04 12:00:00'
		grant user User1 GET /service/service1
		---------------------------------------------------------

		# Create a policy with name "p01" using pdl
		spctl create policy p01 --pdl-command "grant group Administrators list,watch,get expr:c1/default/core/pods/*" --service-name=service1

		# Create a poliy in service service1 using the data in policy.json.
		spctl create policy --json-file ./policy.json --service-name=service1

		# Create a role policy with name "rp01" using pdl
		spctl create rolepolicy rp01 --pdl-command "grant user User1 Role1 on res1" --service-name=service1

		# Create a role poliy in service service1 using the data in rolePolicy.json.
		spctl create rolepolicy --json-file ./rolePolicy.json --service-name=service1
		
		# Create a function "foo", funcUrl , cacheResult, cacheTTL 
		spctl create function foo --func-url=https://a.b.c:3456/funcs/foo --cachable=true --cache-ttl=3600

		# Create a function using function definition json file
		spctl create function --json-file=function.json`
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create (service | policy | rolepolicy | function) (NAME | --json-file JSON_FILENAME) [--pdl-command COMMMAND] [--service-type=TYPE] [--pdl-file=PDL FILE NAME] [--service-name=NAME]",
		Short:   "Create a service | policy | role-policy",
		Example: createExample,
		Run:     createCommandFunc,
	}

	cmd.Flags().StringVarP(&serviceType, "service-type", "t", pms.TypeApplication, "service type, e.g. k8s")
	cmd.Flags().StringVarP(&serviceName, "service-name", "s", "", "service name")
	cmd.Flags().StringVarP(&command, "pdl-command", "c", "", "policy definition language command")
	cmd.Flags().StringVarP(&jsonFileName, "json-file", "f", "", "file that contains policy/role policy/service/function definition in json format")
	cmd.Flags().StringVarP(&pdlFileName, "pdl-file", "l", "", "file that contains policy/role policy definition in policy definition language format")
	cmd.Flags().StringVarP(&funcURL, "func-url", "", "", "URL for the function")
	cmd.Flags().BoolVarP(&funcResultCachable, "cachable", "", false, "whether the function result is cachable")
	cmd.Flags().Int64VarP(&funcResultTTL, "cache-ttl", "", 0, "How many seconds could the function result be kept in cache, 0 means the result could be kept in cache forever")
	return cmd
}

func parsePdlFile(pdlFileName string, serviceName, serviceType string) (*pms.Service, error) {
	f, err := os.Open(pdlFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	service := pms.Service{Name: serviceName, Type: serviceType}
	isRolePolicy := false
	isPolicy := false
	for i, line := range lines {
		line = strings.Trim(line, " ")
		line = strings.Trim(line, "\t")
		if "policies:" == line {
			isPolicy = true
			isRolePolicy = false
		} else if "role-policies:" == line {
			isPolicy = false
			isRolePolicy = true
		} else {
			if isPolicy {
				var policy *pms.Policy
				name := pdlFileName + "_policy_" + strconv.Itoa(i+1)
				policy, _, err := pdl.ParsePolicy(line, name)
				if err == nil {
					service.Policies = append(service.Policies, policy)
				}
			}
			if isRolePolicy {
				var rolePolicy *pms.RolePolicy
				name := pdlFileName + "_role_policy_" + strconv.Itoa(i+1)
				rolePolicy, _, err := pdl.ParseRolePolicy(line, name)
				if err == nil {
					service.RolePolicies = append(service.RolePolicies, rolePolicy)
				}
			}
		}
	}

	return &service, err
}

func createCommandFunc(cmd *cobra.Command, args []string) {
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
	var res string

	switch strings.ToLower(args[0]) {
	case "service":
		var buf []byte
		if len(args) == 1 {
			if jsonFileName == "" {
				cmd.Help()
				return
			}
			buf, err = ioutil.ReadFile(jsonFileName)

		} else if len(args) == 2 {
			serviceName = args[1]
			if serviceName == "" || serviceType == "" {
				cmd.Help()
				return
			}

			if pdlFileName == "" {
				service := pms.Service{Name: serviceName, Type: serviceType}
				buf, err = json.Marshal(service)
			} else {
				var service *pms.Service
				service, err = parsePdlFile(pdlFileName, serviceName, serviceType)
				if err == nil {
					buf, err = json.Marshal(service)
				}
			}

		}
		if err == nil {
			res, err = cli.Post([]string{"service"}, bytes.NewBuffer(buf), "")
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
		if command != "" {
			var buf io.Reader
			var name string
			if len(args) == 2 {
				name = args[1]
			}
			if kind == "policy" {
				_, buf, err = pdl.ParsePolicy(command, name)

			} else {
				_, buf, err = pdl.ParseRolePolicy(command, name)
			}
			if err == nil {
				res, err = cli.Post([]string{"service", serviceName, kind}, buf, "")
			}
		} else {
			if len(args) != 1 || jsonFileName == "" {
				cmd.Help()
				return
			}
			var buf []byte
			buf, err = ioutil.ReadFile(jsonFileName)
			if err == nil {
				res, err = cli.Post([]string{"service", serviceName, kind}, bytes.NewBuffer(buf), "")
			}
		}
	case "function":
		var buf []byte
		if len(args) == 1 {
			if jsonFileName == "" {
				cmd.Help()
				return
			}
			buf, err = ioutil.ReadFile(jsonFileName)

		} else if len(args) == 2 {
			funcName := args[1]
			function := pms.Function{
				Name:           funcName,
				FuncURL:        funcURL,
				ResultCachable: funcResultCachable,
				ResultTTL:      funcResultTTL,
			}
			buf, err = json.Marshal(function)

		}
		if err == nil {
			res, err = cli.Post([]string{"function"}, bytes.NewBuffer(buf), "")
		}

	default:
		cmd.Help()
		return
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("%s created\n%s\n", args[0], res)
	}

}
