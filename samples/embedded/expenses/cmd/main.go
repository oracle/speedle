//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/oracle/speedle/pkg/store/file"
	"github.com/oracle/speedle/samples/embedded/expenses"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <spdl file>.\n", os.Args[0])
		os.Exit(1)
	}
	spdlLoc := os.Args[1]

	handler, err := expenses.Wrap(expenses.ExpenseHTTPHandler, "expenses", spdlLoc)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
