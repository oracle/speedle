package main

import (
	"log"
	"os"
	"os/user"
	"strconv"

	"fmt"
	"github.com/docker/go-plugins-helpers/authorization"
	"speedle/api/authz"
)

const (
	pluginSocket = "/run/docker/plugins/speedle-docker-plugin.sock"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <speedle host> <service name>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	isSecure := "true"
	if os.Args[1] != "a.authz.fun" {
		isSecure = "false"
	}

	properties := map[string]string{
		authz.HOST_PROP:      os.Args[1],
		authz.IS_SECURE_PROP: isSecure,
	}
	plugin, err := newPlugin(properties, os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	// Start service handler on the local sock
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)
	handler := authorization.NewHandler(plugin)
	if err := handler.ServeUnix(pluginSocket, gid); err != nil {
		log.Fatal(err)
	}
}
