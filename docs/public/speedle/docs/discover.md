# Why discover?

Defining authorization policy is always the pain point, especially in system where a large number of resources to be protected.     
We need a way to know *WHAT* actions on *WHAT* resources are protected before defining policies for a system.

# What discover provides?

### 1. Help to discover authorization requests
When *is-allowed* url is replaced with *discover* url, all authorization requests are recorded and "is-allowed=true" is returned at the same time.    
Each authorization request shows *WHO* carries out *WHAT* action on *WHAT* resource.

### 2. Help to generate policies 
Generate policies based on authorization requests recoreded.    
The generated policies allow authorization requests recoreded get "is-allowed= true" result for *is-allowed* call.    
The generated policies could be imported directly.

### 3. Disable authorization without code change
*is-allowed* REST endpoint and *discover* REST endpoint has the same format of request and response. *discover* endpoint always returns "is-allowed=true", at the same time records the authorization request.
When *is-allowed* url is replaced with *discover* url, all authorization check get "is-allowed=true" result. It's a conveneinet way for developers to disable authorization check when diagnose failures.

# How to discover?

### step 1. replace **is-allowed** url with **discover** url in your system
is-allowed REST call is scattered here and there in a system where resources are to be protected. The URL of is-allowed is usually configured somewhere.    
The first step is to replace the url of is-allowed with the url of discover in the configuration, and make the configuration take effect.
```
http://localhost:6734/authz-check/v1/is-allowed ---> http://localhost:6734/authz-check/v1/discover
```
### step 2. [optional] discover authorization requests continuously using CLI(spctl)
```
spctl discover request --last --force --service-name=YOUR_SERVICE_NAME
```
keep the window open to review requests in step 3

### step 3. access resources which are protected in your system
This could trigger a lot of authorization requests be sent out. 

### step 4. generate policies for your service based on authorization requests in step 3.
Using CLI(spctl):
```
spctl discover policy --service-name=YOUR_SERVICE_NAME > service.json
```
### step 5. [optional] import policies generated in step 4 using CLI(spctl)
```
spctl create service --json-file service.json
```
### step 6. [optional] change discover url back to is-allowed url
Don't forget to change back the url when policy is properly created.
    
    
# Discover CLI Reference:    
```
[cyding@test-3 bin]$ ./spctl discover --help
discover request or policy for services

Usage:
  spctl discover (request/policy/reset  | --service-name=NAME | --last | --force | --principal-name=USERNAME) [flags]

Examples:

        # List all request details for all services
        spctl discover request

        # List all request details for the given service
        spctl discover request --service-name="foo"
		
        # List the last request details for service "foo" 
        spctl discover request --last --service-name="foo"
        
        # List the latest request details for service "foo",  doesn't exit until kill it by "Ctrl-C"
        spctl discover request --last --service-name="foo" -f       

        # cleanup all requests
        spctl discover reset

        # clean up the requests for service "foo"
        spctl discover reset --service-name="foo"

        # Generate JSON based policy definition, all users are converted to a role. For example, user Jon visited resourceA. then the following policy will be generated "grant role role_Jon visit resourceA"
        spctl discover policy  --service-name="foo"

        # Generaete JSON based policy definition, only for discover requests triggered by principal which has name 'Jon'
        spctl discover policy --principal-name="Jon" --service-name="foo"

Flags:
  -f, --force                   continuously discover last request
  -h, --help                    help for discover
  -l, --last                    list last request
      --principal-IDD string    principal Identity Domain
      --principal-name string   principal name
      --principal-type string   principal type, could be 'user', 'group','entity'
  -s, --service-name string     service name


```