
# Alpha

Single user

Deploy Action from computer to a remote server

Run from web browser

Dashboard screen:
    List of actions
    
Add trigger:
    Webhook
    Schedule

Enable to create a chani of actions:
    This can be bootstrapped using Webhooks and then with internal routing
    
Store secrets securely

Find sandboxing solution
## Highlevel design

```
 
    |----------|      |----------|        
    |  Engine  |      | Register |  
    |----------|      |----------|

```

#### Packer

on users machine


Deployer

Register

Engine





#### Tasks
- [x] Allow Developer to Deploy the Action:
  
  - [x] Allow to deploy from current working dir or path  
  
- [x] Parse flags

- [x] Allow Developer to Run the Action

- [x] Allow Developer to List all the Actions They Previously Deployed

- [ ] Add named params to cli

- [ ] Validate params

- [ ] Fix logging 

- [ ] Create a server

- [ ] Add Action Schedule

- [x] Add Go support

- [ ] Expose Action as API Endpoint

- [ ] Add Stages/Envs When Deploying an Action 

- [ ] Create Action Skeleton 

- [ ] Log all server actions(deploy, run, trigger, shutdown, restart)

- [ ] Deploy Remotely 