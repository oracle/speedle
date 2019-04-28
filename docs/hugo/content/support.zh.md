+++
draft= false
title = "支持"
description = "Seek answers and ask questions to the Speedle community."
toc = false
tocheading = "h2"
tocsidebar = false
icon = "8. Community.svg"
+++

## FAQs

### Speedle looks awesome! Who is it for?

Speedle is an open source project for access management. You can embed it in your system as long as access control is required, such that you don't have to hard code authorization policies in your code.

### Does the world need one more policy engine?

Well, please spend 15 minutes on [Quick Start](./quick-start) steps. Then you will have the answer. Why not give a try? :-)

### How do I contribute?

Please refer to [contribution guide](./contribute/#code)

### What does the future hold for Speedle? What’s the plan?

We will maintain Speedle as an open source project, fix bugs, add more features...

### How is it different from OPA, Kubernetes RBAC, Istio RBAC?

Long story short: They all target to authorization/access control. Kubernetes RBAC and Istio RBAC focus on the authorization issues in Kubernetes/Istio system respectively, while Speedle and OPA are general purpose authorization solution. There are also a few other general purpose authorization solutions, like Casbin and Ladon. You may try them out yourself and see

1. Is the product user-friendly? These products all offer policy definition language. Try to create a policy and pick up the one you like
2. Is the product fast enough? i.e. how long does it take to evaluate an authorization request?
3. Is the product is scalable enough?
4. Is the supported policy model flexible enough?

Please tell us your choice.

### Why does Speedle have no built-in auth mechanisms?

Speedle focuses on authorization field. It doesn't cover authentication functionality.

### How do I plug in my own backend store?

Please refer to [Pluggable Storage](./docs/store/)

### How do I plug in my own token assertor?

Plase refer to [token assertor](./docs/assertor)

<div class="row" style="padding-top: 50px">
    <div class="col col-2"></div>
    <div class="col col-8 center" style="font-size: 25px">
    <div>
        <p>Can’t see the answer you’re looking for?</p>
        <p>Try our Slack workspace #speedleproject</p>
        <p>
          <a class="button primary started" style="font-size:18px" href="./quick-start">
            <span>Visit Slack channel →</span>
          </a>
        </p>
    </div>
    </div>
    <div class="col col-2 right">
      <a title="slack" href="https://join.slack.com/t/speedleproject/shared_invite/enQtNTUzODM3NDY0ODE2LTg0ODc0NzQ1MjVmM2NiODVmMThkMmVjNmMyODA0ZWJjZjQ3NDc2MjdlMzliN2U4MDRkZjhlYzYzMDEyZTgxMGQ">
        <img class="svg" src="/img/speedle/Slack_RGB.svg" />
      </a>
    </div>
</div>
