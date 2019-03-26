## What is Identity Domain? ##

Identiy domain is a concept from Identity Management. Each identity domain manages users, groups independently. User name and group name are unique in an identity domain.        
In an integrated environment, identities may come from multiple identity domains. To uniquely specify a user/group in such environment, name is not enough, we have to specify identity domain of the user/group. In such environment, principals got after authentication has identity domain specified.

## Identity Domain Support in Speedle ##
### Policy Definition ###
#### 1. principal with identity domain in SPDL ####  

use 'from IDENTITY_DOMAIN_NAME' to indicate identity domain of a principal in SPDL  
```
from IDENTITY_DOMAIN_NAME
```  
sample role policy in SPDL:
```
grant user USER_NAME from IDENTITY_DOMAIN_NAME ROLENAME
grant group GROUP_NAME from IDENTITY_DOAMIN_NAME ROLENAME
```
#### 2. principal with identity domain in json presentation of policy/role policy #### 

In json presentation of a policy or role policy, principal is encoded as: 
```
IDENTITY_DOMAIN_NAME:PRINCIPAL_TYPE:PRINCIPAL_NAME  
```     
	
see principals in the sample json definition of policy    
```
{                    
	"name": "allowViewersView",
	"effect": "grant",
	"permissions": [
		{
			"resourceExpression": "/service/*",
			"actions": ["GET"]
		}
	],
	"principals": [	"identityDomain1:group:Viewers"	],
	"conditions": [	"request_time < '2017-09-04 12:00:00'" ]
}
```
### Authorization check ###
#### principal with identity domain in ADS request ####  

In an ADS request, principal is defined as follows, identity domain could be set via "IDD" property:   
```
type Principal struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
	IDD  string `json:"idd,omitempty"`
}
```
Following is a sample json ADS request, in which the principal has identity domain info:    
```

{
 "subject": {"principals":[{"idd":"identityDomain_1", "type":"user", "name":"cyding"}]},
 "serviceName": "service1",
 "resource":"laptop",
 "action": "access"
}
```









