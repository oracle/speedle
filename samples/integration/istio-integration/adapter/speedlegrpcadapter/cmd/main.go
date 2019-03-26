package main

import (
	"fmt"
	"os"

	"istio.io/istio/mixer/adapter/speedlegrpcadapter"
)

func main() {

	addr := ""
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	s, err := speedlegrpcadapter.NewGrpcServer(addr)
	if err != nil {
		fmt.Printf("unable to start server: %v", err)
		os.Exit(-1)
	}

	s.Run()
	s.Wait()
}
