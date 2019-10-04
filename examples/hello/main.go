package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Handler struct {}

func (handler *Handler) Action(data []byte, result *[]byte) error {

	var params map[string]interface{}
	err := json.Unmarshal(data, &params)
	if err != nil {
		return err
	}

	resp := []byte(fmt.Sprintf("Hello,%v", params["name"]))

	result = &resp

	return nil
}

func Invoke(port int) {
	h := &Handler{}

	err := rpc.Register(h)
	if err != nil {
		log.Fatal()
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))

	log.Println("listening on", port)
	err = http.Serve(l, nil)
	if err != nil {

	}
}

func main() {

	var port int
	flag.IntVar(&port, "port", 6060, "")

	flag.Parse()
	Invoke(port)
}


