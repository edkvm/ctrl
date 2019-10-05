package main

import (
	"fmt"
	"github.com/edkvm/ctrl/gowrapper"
)


type Params struct{
	Name string `json:"name"`
}

func Hello(params Params) (string, error) {
	return fmt.Sprintf("Hello, %v", params.Name), nil
}

func main() {
	gowrapper.Start(Hello)
}


