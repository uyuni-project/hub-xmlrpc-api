package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/divan/gorilla-xmlrpc/xml"
	"github.com/gorilla/rpc"
)

type Auth struct{}

func (h *Auth) AttachToServer(r *http.Request, args *struct{ HubSessionKey, ServerURL, Username, Password string }, reply *struct{ Message string }) error {
	//TODO: check the hubToken
	response, _ := XMLRPCCall(args.ServerURL, "auth.login", &struct{ Username, Password string }{args.Username, args.Password})
	reply.Message = response.Data.(string)
	return nil
}

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Message string }) error {
	loginToken := "1234"
	reply.Message = loginToken
	return nil
}

type DefaultService struct{}

type DefaultCallParams struct {
	HubKey string
	Elems  [][]interface{}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *DefaultCallParams, reply *struct{ Data map[string]interface{} }) error {
	//TODO: check the hubToken

	endpoints := []string{"http://192.168.122.76/rpc/api"}

	method, _ := xml.NewCodec().NewRequest(r).Method()

	responses := make(map[string]interface{})

	for i, url := range endpoints {

		arguments := args.Elems[i]

		response, err := XMLRPCCall(url, method, &arguments)
		if err != nil {
			log.Fatal(err)
		}
		responses[url] = response.Data
		log.Printf("Response: %s\n", response.Data)
	}

	reply.Data = responses
	return nil
}

func XMLRPCCall(url string, method string, args interface{}) (reply struct{ Data interface{} }, err error) {
	buf, _ := xml.EncodeClientRequest(method, args)

	resp, err := http.Post(url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = xml.DecodeClientResponse_(resp.Body, &reply)
	return
}

func main() {
	RPC := rpc.NewServer()
	xmlrpcCodec := xml.NewCodec()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	RPC.RegisterService(new(Auth), "")
	RPC.RegisterDefaultService(new(DefaultService), "DefaultService.DefaultMethod")

	http.Handle("/RPC2", RPC)

	log.Println("Starting XML-RPC server on localhost:8000/RPC2")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
