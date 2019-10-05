### The Problem

Developers spend a lot of their time on routine and mundane tasks  
like migration, data entry, data enrichment, data processing, connecting between
various apis and ect. 
 
Those task are usually not directly business related and time spent on them   
is minimal, though usually wastful.   
 
The tasks are considered a disposable solution to a very specific problem they   
are not designed to be used again. Unfortunately it is usually not the case and    
developers find themselves implementing similar task within the same organization.   

Developing these tasks is annoying, running or repeating them in a later  
stage is a major hassle.



### Current Solutions 

Most existing solutions look at the problem from a client point of view, instead of a developer.   
The tools eventual converage to be more usable through a UI. This is not a bad idea per say,
but it lacks the focus.

Another fallacy is that current solutions are advertised as a replacment to production work
loads, which is far from reality.

#### Related Ideas
- StdLib
- Zapier
- IFFFT
- TriggerHappy(OSS)
- Huginn(OSS)


### Proposed solution

A Distributed Personalized CI

- Open Platform
- Personal CI?
- Developer Runtime?

An easy way to build tools?
Project Automation


Must Have:
* Observability
* Retry
* Easy Modifications for parameters
* Deploy Stages
* Security
* Hooks
 

Should Have:
* Consider using git interface
* Build on the server
* Access for the web with an endpoint


### Assumptions
Developers like to write code.

Learning how to integrate or interact with closed platforms like zapier  
is considered a waste of time(Solution: open platform, simple interface).

Developers write a lot of open source code which is mostly useless  
only big and prominent projects get picked up.

Allowing developers to write small, less than 100 LOC, programs that can  
be used right away by third party collaborators is a big value proposition  
for developers.

Developing the correct interface is crucial, multiple interface should be studied  
AWS Lambda, Google Function, Azure Functions.

None of the commercial solution can be used on premise

Performance is not the big issues








