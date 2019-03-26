# Authorization diagnosis

This feature is used to diagnose the evaluation process of authorization. When a user isn't allowed to operate on a resource, then this feature is useful to find out the reason, i.e., what's the policy which denies the operation on the resource.

## Format of diagnosis response

The format of the diagnosis response is as below.

```go
type EvaluationDebugResponse struct {
	Allowed       bool                    `json:"allowed"`
	Reason        string                  `json:"reason"`
	RequestCtx    *RequestContext         `json:"requestContext,omitempty"`
	Attributes    map[string]interface{}  `json:"attributes,omitempty"`
	GrantedRoles  []string                `json:"grantedRoles,omitempty"`
	RolePolicies  []*EvaluatedRolePolicy  `json:"rolePolicies,omitempty"`
	Policies      []*EvaluatedPolicy      `json:"policies,omitempty"`
}

type EvaluatedPolicy struct {
	Status      string              `json:"status,omitempty"`
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Effect      string              `json:"effect,omitempty"`
	Permissions []Permission        `json:"permissions,omitempty"`
	Principals  []string            `json:"principals,omitempty"`
	Condition   *EvaluatedCondition `json:"condition,omitempty"`
}

type EvaluatedRolePolicy struct {
	Status              string              `json:"status,omitempty"`
	ID                  string              `json:"id,omitempty"`
	Name                string              `json:"name,omitempty"`
	Effect              string              `json:"effect,omitempty"`
	Roles               []string            `json:"roles,omitempty"`
	Principals          []string            `json:"principals,omitempty"`
	Resources           []string            `json:"resources,omitempty"`
	ResourceExpressions []string            `json:"resourceExpression,omitempty"`
	Condition           *EvaluatedCondition `json:"condition,omitempty"`
}

type EvaluatedCondition struct {
	ConditionExpression  string  `json:"conditionExpression,omitempty"`
	EvaluationResult     string  `json:"evaluationResult,omitempty"`
}
```

## Example(s)

Note that both the attributes provided by the request and the build-in attributes are included in the response.

The field "status" in each policy or rolePolicy has three possible values, which are "takeEffect", "conditionFailed", or "ignored",

-   takeEffect

"takeEffect" means that the policy or rolePolicy is matched and evaluated.

-   conditioFailed

"conditionFailed" means that the policy or rolePolicy matches the service name, subject, resource and action(Note that action only applies to policy), but the evaluation result of the condition is false.

-   ignored

"ignored" means that the evaluation process of authorization is already done, so no need to evaluate the policy any more.

The following is an example,

```go
{
  "Allowed": "true",
  "requestContext": {
    "subject": {
      "user": "user1",
      "groups": null,
      "attributes": null
    },
    "serviceName": "srv1",
    "resource": "res1",
    "action": "read",
    "attributes": null,
    "token": null
  },
  "attributes": {
    "request_action": "read",
    "request_day": 23,
    "request_groups": null,
    "request_month": "November",
    "request_resource": "res1",
    "request_time": 1511406017,
    "request_user": "user1",
    "request_weekday": "Thursday",
    "request_year": 2017
  },
  "grantedRoles": [
    "role1"
  ],
  "rolePolicies": [
    {
      "status": "takeEffect",
      "id": "c8087db3-60cf-4dad-aa9d-033eb6da0b15",
      "name": "rp01",
      "effect": "grant",
      "roles": [
        "role1"
      ],
      "principals": [
        "user:user1",
        "user:user2"
      ],
      "resources": [
        "res1",
        "res2"
      ],
      "condition": {

      }
    }
  ],
  "policies": [
    {
      "status": "takeEffect",
      "id": "f56b494f-dd6b-42af-962e-a109c890b7a0",
      "name": "p01",
      "effect": "grant",
      "permissions": [
        {
          "resource": "res1",
          "actions": [
            "list",
            "read",
            "write"
          ]
        },
        {
          "resource": "res2",
          "actions": [
            "list"
          ]
        }
      ],
      "principals": [
        "user:user1",
        "user:user2"
      ],
      "condition": {
        "conditionExpression": "request_year ==2017",
        "evaluationResult": "true"
      }
    }
  ]
}
```
