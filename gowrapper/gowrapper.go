package gowrapper

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"reflect"
)

type Action struct {
	handler interface{}
}

func NewAction(handler interface{}) *Action{
	return &Action{
		handler: handler,
	}
}

type InvokeReq struct {
	ID      string `json:"id"`
	Payload []byte `json:"payload"`
}

type InvokeRes struct {
	ID      string `json:"id"`
	Payload []byte `json:"payload"`
}

func (ac *Action) Invoke(input []byte, out *[]byte) error {

	var req InvokeReq

	err := json.Unmarshal(input, &req)
	if err != nil {
		log.Println(err)
		return err
	}

	// Grab Handler
	handlerType := reflect.TypeOf(ac.handler)

	// Grab first Param
	paramType := handlerType.In(0)

	// Create the param
	ptrValue := reflect.New(paramType)

	pValue := ptrValue.Interface()

	err = json.Unmarshal(req.Payload, pValue)
	if err != nil {
		return err
	}


	in := []reflect.Value{ptrValue.Elem()}
	m := reflect.ValueOf(ac.handler)

	actionRes := m.Call(in)

	encoded, err := json.Marshal(actionRes[0].Interface())

	resp := &InvokeRes{
		Payload: encoded,
	}


	encResp, err := json.Marshal(resp)
	if err != nil {

	}

	*out = encResp
	return nil
}



func Start(handler interface{}) {
	fd := os.Getenv("CTRL_INT_SOCKET")
	h := NewAction(handler)

	err := rpc.Register(h)
	if err != nil {
		log.Fatal(err)
	}
	rpc.HandleHTTP()

	l, err := net.Listen("unix", fd)
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(l, nil)
	if err != nil {

	}
}

