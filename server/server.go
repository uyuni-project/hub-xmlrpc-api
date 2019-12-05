package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/chiaradiamarcelo/hub_xmlrpc_api/client"
	"github.com/chiaradiamarcelo/hub_xmlrpc_api/config"
	"github.com/chiaradiamarcelo/hub_xmlrpc_api/session"
	"github.com/gorilla/rpc"
)

var conf = config.New()
var apiSession = session.New()

type DefaultService struct{}

type DefaultCallArgs struct {
	HubKey     string
	ServerIds  []int64
	ServerArgs [][]interface{}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *DefaultCallArgs, reply *struct{ Data []interface{} }) error {
	if apiSession.IsHubSessionValid(args.HubKey) {
		method, err := NewCodec().NewRequest(r).Method()
		if err != nil {
			log.Println("Call error: %v", err)
		}
		//TODO: check args.ServerArgs lists have the same size
		serverArgsByURL := make(map[string][]interface{})

		for i, serverID := range args.ServerIds {
			out := make([]interface{}, len(args.ServerArgs)+1)

			for j, serverArgs := range args.ServerArgs {
				out[j+1] = serverArgs[i]
			}
			url, sessionKey := apiSession.GetServerSessionInfoByServerID(args.HubKey, serverID)
			out[0] = sessionKey
			serverArgsByURL[url] = out
		}
		reply.Data = multicastCall(method, serverArgsByURL)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func multicastCall(method string, serverArgsByURL map[string][]interface{}) []interface{} {
	responses := make([]interface{}, len(serverArgsByURL))

	var wg sync.WaitGroup
	wg.Add(len(serverArgsByURL))

	i := 0
	for url, args := range serverArgsByURL {
		go func(url string, args []interface{}, i int) {
			defer wg.Done()
			response, err := executeXMLRPCCall(url, method, args)
			if err != nil {
				log.Println("Call error: %v", err)
			}
			responses[i] = response
			log.Printf("Response: %s\n", response)
		}(url, args, i)
		i++
	}
	wg.Wait()
	return responses
}

func executeXMLRPCCall(url string, method string, args []interface{}) (reply interface{}, err error) {
	client, err := client.GetClientWithTimeout(url, conf.ConnectTimeout, conf.ReadWriteTimeout)
	if err != nil {
		return
	}
	defer client.Close()
	err = client.Call(method, args, &reply)
	return reply, err
}

func InitServer() {
	xmlrpcCodec := NewCodec()
	xmlrpcCodec.RegisterMethod("Hub.Login")
	xmlrpcCodec.RegisterMethod("Hub.AttachToServers")
	xmlrpcCodec.RegisterMethod("Hub.ListServerIds")
	xmlrpcCodec.RegisterDefaultMethod("DefaultService.DefaultMethod")

	RPC := rpc.NewServer()
	RPC.RegisterCodec(xmlrpcCodec, "text/xml")
	RPC.RegisterService(new(Hub), "")
	RPC.RegisterService(new(DefaultService), "")

	http.Handle("/RPC2", RPC)

	log.Println("Starting XML-RPC server on localhost:8000/RPC2")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
