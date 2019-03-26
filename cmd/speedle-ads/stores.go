//Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

// In this file, we link every data store implmention with a side-effect import (using a blank import name). You can add your own store here too.
// If you want to use speedle as in-process mode, you can copy this stores.go to your own package and modify the package name to your own package name.

package main

import (
	_ "github.com/oracle/speedle/pkg/store/etcd"
	_ "github.com/oracle/speedle/pkg/store/file"
)
