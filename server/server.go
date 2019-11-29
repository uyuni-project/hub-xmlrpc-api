package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/chiaradiamarcelo/hub_xmlrpc_api/client"
	"github.com/gorilla/rpc"
)

type Auth struct{}

func (h *Auth) AttachToServer(r *http.Request, args *struct{ HubSessionKey, ServerURL, Username, Password string }, reply *struct{ Message string }) error {
	//TODO: check the hubToken
	response, _ := executeXMLRPCCall(args.ServerURL, "auth.login", []interface{}{args.Username, args.Password})
	reply.Message = response.(string)
	return nil
}

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	reply.Data = "1234"
	return nil
}

type DefaultService struct{}

type DefaultCallParams struct {
	HubKey string
	Elems  [][]interface{}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *DefaultCallParams, reply *struct{ Data map[string]interface{} }) error {
	//TODO: check the hubToken

	endpoints := []string{"http://192.168.122.76/rpc/api", "http://192.168.122.2/rpc/api"}

	method, _ := NewCodec().NewRequest(r).Method()

	responses := make(map[string]interface{})
	//Execute the calls concurrently but wait before we get the response from all the servers.
	var wg sync.WaitGroup
	wg.Add(len(endpoints))

	for i, url := range endpoints {
		go func(url string, args []interface{}) {
			defer wg.Done()
			response, err := executeXMLRPCCall(url, method, args)
			if err != nil {
				log.Println("Call error: %v", err)
			}
			responses[url] = response
			log.Printf("Response: %s\n", response)
		}(url, args.Elems[i])
	}
	wg.Wait()
	reply.Data = responses
	return nil
}

func executeXMLRPCCall(url string, method string, args []interface{}) (reply interface{}, err error) {
	client, err := client.GetClientWithTimeout(url, 2, 5)
	if err != nil {
		return
	}
	defer client.Close()

	err = client.Call(method, args, &reply)

	return reply, err
}

func InitServer() {
	xmlrpcCodec := NewCodec()
	xmlrpcCodec.RegisterMethod("Auth.Login")
	xmlrpcCodec.RegisterMethod("Auth.AttachToServer")
	xmlrpcCodec.RegisterDefaultMethod("DefaultService.DefaultMethod")

	RPC := rpc.NewServer()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	RPC.RegisterService(new(Auth), "")
	RPC.RegisterService(new(DefaultService), "")

	http.Handle("/RPC2", RPC)

	log.Println("Starting XML-RPC server on localhost:8000/RPC2")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
