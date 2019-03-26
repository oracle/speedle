package client

import "speedle/api/authz"

// ADSClient is a client interface for ADS service
type ADSClient interface {
	IsAllowed(authz.RequestContext) (bool, error)
}
