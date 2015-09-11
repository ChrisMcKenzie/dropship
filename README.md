# Dropship

Dropship is a service that sits on the edge of your network and the internet
and lets you perform automated deployments via `Github Deployments`. It 
accomplishes this by first retrieving a `dropship.yml` file from your repo
which defines a few things such as task that it will run on the also defined servers.

```
+---------+             +--------+            +----------+         +-------------+
| Tooling |             | GitHub |            | Dropship |         | Your Server |
+---------+             +--------+            +----------+         +-------------+
     |                      |                       |                     |
     |  Create Deployment   |                       |                     |
     |--------------------->|                       |                     |
     |                      |                       |                     |
     |  Deployment Created  |                       |                     |
     |<---------------------|                       |                     |
     |                      |                       |                     |
     |                      |   Deployment Event    |                     |
     |                      |---------------------->|                     |
     |                      |                       |     SSH+Deploys     |
     |                      |                       |-------------------->|
     |                      |                       |                     |
     |                      |   Deployment Status   |                     |
     |                      |<----------------------|                     |
     |                      |                       |                     |
     |                      |                       |   Deploy Completed  |
     |                      |                       |<--------------------|
     |                      |                       |                     |
     |                      |   Deployment Status   |                     |
     |                      |<----------------------|                     |
     |                      |                       |                     |
```

## Roadmap
- [X] Deployment via ssh on server.
- [ ] Environment specific deployments [#1](https://github.com/ChrisMcKenzie/dropship/issues/1)
- [ ] Multiple task definitions
- [ ] Deployment Logging and UI
- [ ] Setup UI
- [ ] Multiple users can view project

