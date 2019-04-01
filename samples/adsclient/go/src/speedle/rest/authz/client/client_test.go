package client

import (
	"log"
	"speedle/api/authz"
	"testing"
)

var connectionProperties = map[string]string{}

func TestIsAllowed(t *testing.T) {
	client, err := New(connectionProperties)
	if err != nil {
		log.Printf("Error in creating a new client due to error %v.\n", err)
		t.FailNow()
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
		t.FailNow()
	}
	log.Printf("Evaluation result: %v.", allowed)
	if !allowed {
		log.Printf("Wrong result %v.", allowed)
	}
}
