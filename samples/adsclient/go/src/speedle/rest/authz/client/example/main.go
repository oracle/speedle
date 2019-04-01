package main

import (
	"log"
	"os"
	"speedle/api/authz"
	"speedle/rest/authz/client"
)

func main() {
	connectionProperties := map[string]string{
		authz.IS_SECURE_PROP: "false",
	}

	client, err := client.New(connectionProperties)
	if err != nil {
		log.Printf("Error in creating a new client due to error %v.\n", err)
		os.Exit(1)
	}

	context := authz.RequestContext{
		Subject: &authz.Subject{
			Principals: []*authz.Principal{
				{
					Type: "user",
					Name: "alan.cao",
				},
			},
		},
		ServiceName: "acao",
		Resource:    "/home/acao/Downloads",
		Action:      "read",
	}
	allowed, err := client.IsAllowed(context)
	if err != nil {
		log.Printf("Error in calling IsAllowed due to error %v.\n", err)
		os.Exit(1)
	}
	log.Printf("Evaluation result: %v.", allowed)
}
