//Copyright (c) 2018, Oracle and/or its affiliates. All rights reserved.
//Licensed under the Universal Permissive License (UPL) Version 1.0 as shown at http://oss.oracle.com/licenses/upl.

package subjectutils

import (
	"fmt"

	adsapi "github.com/oracle/speedle/api/ads"
)

// EncodePrincipal encodes prinicpal object to string
// Form: [idd=<IDD>:]<Type>:<Name>
func EncodePrincipal(principal *adsapi.Principal) string {
	if len(principal.IDD) != 0 {
		return fmt.Sprintf("idd=%s:%s:%s", principal.IDD, principal.Type, principal.Name)
	}
	return fmt.Sprintf("%s:%s", principal.Type, principal.Name)
}
