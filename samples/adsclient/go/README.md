# Speedle Golang Client

This is a Golang Client for Speedle ADS to simplify ADS calls.

# Interface
```go
package client

import "speedle/api/authz"

// ADSClient is a client interface for ADS service
type ADSClient interface {
    IsAllowed(authz.RequestContext) (bool, error)
}
```

# How to construct a REST client instance
```go
client, err := client.New("127.0.0.1", false)
```

Construct a client with

* Speedle host: ```a.authz.fun```
* Insecure? ```false```

# Example
The file [main.go](src/speedle/rest/authz/client/example/main.go) showes how to call ADS promatically.

## How to run Example
```sh
go run sphnix/rest/authz/client/example/main.go
```

# Q&A
## How to connect ADS with this client behind a HTTP proxy.
The Client reads following envionment variables for proxy setting:
* HTTPS_PROXY: The PROXY endpoint
* NO_PROXY: Hosts that should not be connected with the proxy.

More details please refer the golang document https://golang.org/pkg/net/http/#ProxyFromEnvironment


