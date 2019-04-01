//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package command

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/oracle/speedle/cmd/spctl/client"
	"github.com/oracle/speedle/pkg/svcs/pmsrest"
	"github.com/spf13/cobra"
)

var (
	last                                       bool
	force                                      bool
	principalType, principalName, principalIDD string
)

var (
	discoverExample = `
        # List all request details for all services
        spctl discover request

        # List all request details for the given service
        spctl discover request --service-name="foo"
		
        # List the last request details for service "foo" 
        spctl discover request --last --service-name="foo"
        
        # List the latest request details for service "foo", doesn't exit until you kill it using "Ctrl-C"
        spctl discover request --last --service-name="foo" -f       

        # cleanup all requests
        spctl discover reset

        # clean up the requests for service "foo"
        spctl discover reset --service-name="foo"

        # Generate JSON based policy definition, all users are converted to a role. For example, user Jon visited resourceA. Then the following policy is generated "grant role role_Jon visit resourceA"
        spctl discover policy  --service-name="foo"

        # Generate JSON based policy definition, only for discover requests triggered by principal which has name 'Jon'
        spctl discover policy --principal-name="Jon" --service-name="foo"`
)

func NewDiscoverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "discover (request/policy/reset  | --service-name=NAME | --last | --force | --principal-name=USERNAME)",
		Short:   "discover request or policy for services ",
		Example: discoverExample,
		Run:     discoverCommandFunc,
	}

	cmd.Flags().BoolVarP(&last, "last", "l", false, "list last request")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "continuously discover last request")
	cmd.Flags().StringVarP(&serviceName, "service-name", "s", "", "service name")
	cmd.Flags().StringVarP(&principalType, "principal-type", "", "", "principal type, could be 'user', 'group','entity'")
	cmd.Flags().StringVarP(&principalName, "principal-name", "", "", "principal name")
	cmd.Flags().StringVarP(&principalIDD, "principal-IDD", "", "", "principal Identity Domain")
	return cmd
}

func discoverCommandFunc(cmd *cobra.Command, args []string) {
	hc, err := httpClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	cli := &client.Client{PMSEndpoint: globalFlags.PMSEndpoint, HTTPClient: hc}
	var res []byte
	var output []byte
	if len(args) != 1 {
		cmd.Help()
		return
	}
	switch strings.ToLower(args[0]) {
	case "request":
		if last {
			if force {
				v := url.Values{}
				v.Add("last", "true")
				res, err = cli.Get([]string{"discover-request", serviceName}, v, "")
				if err == nil {
					var revision int64
					var response pmsrest.GetDiscoverRequestsResponse
					if json.Unmarshal(res, &response) == nil {
						revision = response.Revision
						for _, request := range response.Requests {
							output, _ = json.MarshalIndent(&request, "", strings.Repeat(" ", 4))
							fmt.Println(string(output))
						}
					}
					for {
						v = url.Values{}
						v.Add("revision", strconv.FormatInt(revision, 10))
						res, err = cli.Get([]string{"discover-request", serviceName}, v, "")
						if err == nil {
							if json.Unmarshal(res, &response) == nil {
								revision = response.Revision
								for _, request := range response.Requests {
									output, _ = json.MarshalIndent(&request, "", strings.Repeat(" ", 4))
									fmt.Println(string(output))
								}
							}
						} else {
							break
						}
						time.Sleep(1 * time.Second)
					}
				}
			} else {
				v := url.Values{}
				v.Add("last", "true")
				res, err = cli.Get([]string{"discover-request", serviceName}, v, "")
				if err == nil {
					var response pmsrest.GetDiscoverRequestsResponse
					if json.Unmarshal(res, &response) == nil {
						for _, request := range response.Requests {
							output, _ = json.MarshalIndent(&request, "", strings.Repeat(" ", 4))
							fmt.Println(string(output))
						}
					}
				}
			}
		} else {
			// spxctl discover request --all [--service-name="foo"]
			res, err = cli.Get([]string{"discover-request", serviceName}, nil, "")
			if err == nil {
				var response pmsrest.GetDiscoverRequestsResponse
				err = json.Unmarshal(res, &response)
				if err == nil {
					output, err = json.MarshalIndent(response.Requests, "", strings.Repeat(" ", 4))
					fmt.Println(string(output))
				}
			}
		}
	case "policy":
		if serviceName == "" {
			fmt.Println("pls specify service name by --service-name=NAME")
			return
		} else {
			v := url.Values{}
			if len(principalType) > 0 {
				v.Add("principalType", principalType)
			}
			if len(principalName) > 0 {
				v.Add("principalName", principalName)
			}
			if len(principalIDD) > 0 {
				v.Add("principalIDD", principalIDD)
			}
			res, err = cli.Get([]string{"discover-policy", serviceName}, v, "")
			if err == nil {
				var response pmsrest.GetDiscoverPoliciesResponse
				err = json.Unmarshal(res, &response)
				if err == nil {
					if len(response.Services) > 0 {
						output, _ = json.MarshalIndent(response.Services[0], "", strings.Repeat(" ", 4))
						fmt.Println(string(output))
					} else {
						fmt.Println("no policy discovered for the service")
					}
				} else {
					fmt.Println("fail to unmarshal response,", err)
				}
			} else {
				fmt.Println("fail to discover policy,", err)
			}
		}

	case "reset":
		if serviceName == "" {
			// spxctl reset service --service-name="foo"
			err = cli.Delete([]string{"discover-request"}, "")
			if err == nil {
				fmt.Printf("All requests are deleted.\n")
			}
		} else {
			err = cli.Delete([]string{"discover-request", serviceName}, "")
			if err == nil {
				fmt.Printf("Requests for service(%s) are deleted.\n", serviceName)
			}
		}
	default:
		cmd.Help()
		return
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
