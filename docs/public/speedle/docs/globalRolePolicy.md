## Background ##
In Speedle, role policies and policies are defined and checked in a scope, which is called 'service'. 
For example, when creating a policy, you need to specify in which 'service' the policy is to be created; 
When calling Authorization Decision Service API, a 'service' name should be specified, 
so role policies and policies in that 'service' scope be evaluated.    
### Global Role Policies ###
In some cases, for example, a big system with many sub-systems, each sub-system has its own specific authorization requirement, so 
each sub-system defines role policies and policies in its own 'service' scope.    
But there are also something, like role policies, in common among all sub-systems.    
To avoid creating same role policies in each sub-system, we support global role policies, which are shared among all 'services'.

## Global Service && Global Role Policies ##
### Global Service ###
1. It is a special service named 'global'.
2. It could be created and deleted by customer, just like a normal service.
3. Only role policies could be created in 'global' service. No policies could be created in 'global' service. 
4. User can call ADS API in the service, same as any other service.

### Global Role Policies ###
1. It is the role policies created in 'global' service.
2. Global role policies take effect in policy evaluation in any normal service, just as if those role policies were defined in that service.

## Usage of global role policies ##
### Create global role policies ###
1. create a service named 'global'
```
spctl create service global
```
2. create role policies in 'global' service
```
spctl create rolepolicy -c "grant user Emma AdminRole" --service-name=global
```
### Create policies in a normal service ###
```
spctl create policy -c "grant role AdminRole borrow books" --service-name=library
```
### Call is-allowed API in the normal service scope ###
```
curl -X POST  http://localhost:6734/authz-check/v1/is-allowed \
-d @- << EOF
{
 "subject": {"principals":[{"type":"user", "name":"Emma"}]},
 "serviceName": "library",
 "resource":"books",
 "action": "borrow"
}
EOF
```
allowed = true is returned in this case, since the role policy defined in 'global' service takes effect.
