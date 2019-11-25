package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/divan/gorilla-xmlrpc/xml"
	"github.com/gorilla/rpc"
)

type HelloService struct{}

func (h *HelloService) Say(r *http.Request, args *struct{ Who string }, reply *struct{ Message string }) error {
	log.Println("Say", args.Who)
	reply.Message = "Hello, " + args.Who + "!"
	return nil
}

type DefaultService struct{}

func (h *DefaultService) DefaultMethod(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Message string }) error {
	log.Println("Default", args.Username)

	method := ExtractMethod(r)
	endpoints := []string{"http://192.168.122.76/rpc/api"}

	for _, url := range endpoints {
		response, err := XMLRPCCall(url, method, args)
		if err != nil {
			log.Fatal(err)
		}

		reply.Message = response.Message
		log.Printf("Response: %s\n", response.Message)
	}

	//reply.Message = response
	return nil
}

func ExtractMethod(r *http.Request) string {
	/*xmlrpcCodec := xml.NewCodec()

	buf, _ := xml.EncodeClientRequest(method, &args)

	resp, err := http.Post(url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = xml.DecodeClientResponse(resp.Body, &reply)*/
	return "auth.login"
}

func XMLRPCCall(url string, method string, args interface{}) (reply struct{ Message string }, err error) {
	buf, _ := xml.EncodeClientRequest(method, args)

	resp, err := http.Post(url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = xml.DecodeClientResponse(resp.Body, &reply)
	return
}

func main() {
	RPC := rpc.NewServer()
	xmlrpcCodec := xml.NewCodec()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	RPC.RegisterService(new(HelloService), "")
	RPC.RegisterDefaultService(new(DefaultService), "DefaultService.DefaultMethod")

	http.Handle("/RPC2", RPC)

	log.Println("Starting XML-RPC server on localhost:1234/RPC2")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
