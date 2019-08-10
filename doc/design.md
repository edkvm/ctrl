## Design


#### Action

Data Struct
```go
    {
    	ID string
    	Name string
    	Stack
    	Path
    	Config map[string]string
    }
```


# Questions | Assumptions | Ideas

Support Node 12 and Go

Can only run safe verified programs 

Should Constrain resources

Currently there is a tech challnege to find Sandbox solution  
option to consider: Docker, WebAssembly, NaCl

Easiest option is Docker just run a container and mount a directory

Initial version the build should be done at the users computer

Later we can use the Git model and build on the server

Review heroku build-pack

Use JSON-RPC as possible spec for functions 