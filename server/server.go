package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/chiaradiamarcelo/hub_xmlrpc_api/client"
	"github.com/chiaradiamarcelo/hub_xmlrpc_api/config"
	"github.com/chiaradiamarcelo/hub_xmlrpc_api/session"
	"github.com/gorilla/rpc"
)

var conf = config.New()
var apiSession = session.New()

func isHubSessionValid(in string) bool {
	//TODO: we should check through the SUMA api
	return in == apiSession.GetHubSessionKey()
}

func getServerURLFromServerID(serverID int64) string {
	//TODO:
	return conf.ServerURLByServerID[strconv.FormatInt(serverID, 10)]
}

type DefaultService struct{}

type DefaultCallArgs struct {
	HubKey     string
	ServerArgs [][]interface{}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *DefaultCallArgs, reply *struct{ Data map[string]interface{} }) error {
	if isHubSessionValid(args.HubKey) {
		method, err := NewCodec().NewRequest(r).Method()
		if err != nil {
			log.Println("Call error: %v", err)
		}

		serverArgsByURL := make(map[string][]interface{})

		for _, args := range args.ServerArgs {
			url := apiSession.GetServerURLbyServerKey(args[0].(string))
			serverArgsByURL[url] = args
		}
		reply.Data = multicastCall(method, serverArgsByURL)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func multicastCall(method string, serverArgsByURL map[string][]interface{}) map[string]interface{} {
	responses := make(map[string]interface{})
	//Execute the calls concurrently and wait until we get the response from all the servers.
	var wg sync.WaitGroup

	wg.Add(len(serverArgsByURL))

	for url, args := range serverArgsByURL {
		go func(url string, args []interface{}) {
			defer wg.Done()
			response, err := executeXMLRPCCall(url, method, args)
			if err != nil {
				log.Println("Call error: %v", err)
			}
			responses[url] = response
			log.Printf("Response: %s\n", response)
		}(url, args)
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
