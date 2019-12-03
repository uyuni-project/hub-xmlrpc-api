package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/chiaradiamarcelo/hub_xmlrpc_api/client"
	"github.com/gorilla/rpc"
)

type Auth struct{}

const HUB_SUMA_API_URL = "http://192.168.122.76/rpc/api"

//TODO: session
var hubSessionKey = ""
var username = ""
var pass = ""
var userServerUrlByKey = make(map[string]interface{})

//FLAGS
var autoRelayMode = true
var autoConnectMode = true

var serverEndpoints = []string{"http://192.168.122.76/rpc/api", "http://192.168.122.2/rpc/api"}

func isHubSessionValid(in string) bool {
	//TODO:
	return in == hubSessionKey
}

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data map[string]interface{} }) error {
	sessionkeys := make(map[string]interface{})
	if autoRelayMode {
		//save credentials in session
		username = args.Username
		pass = args.Password

		if autoConnectMode {
			serverKeys, err := loginIntoUserSystems(username, pass)
			if err != nil {
				log.Println("Call error: %v", err)
			}
			sessionkeys["serverKeys"] = serverKeys
		}
	}
	response, _ := executeXMLRPCCall(HUB_SUMA_API_URL, "auth.login", []interface{}{args.Username, args.Password})
	//save in session
	hubSessionKey = response.(string)
	//update respnse
	sessionkeys["hubSessionKey"] = response.(string)
	reply.Data = sessionkeys
	return nil
}

func loginIntoUserSystems(username, pass string) (map[string]interface{}, error) {
	//get user servers
	_, err := executeXMLRPCCall(HUB_SUMA_API_URL, "system.listUserSystems", []interface{}{username, pass})
	if err != nil {
		return nil, err
	}

	//TODO:
	serverCredentials := make([][]interface{}, 1)

	serverCredentials[0] = make([]interface{}, 2)
	serverCredentials[0][0] = struct{ Username, Password string }{"admin", "admin"}
	serverCredentials[0][1] = struct{ Username, Password string }{"admin", "admin"}

	loginResponses := multicastCall("auth.login", serverCredentials)

	for url, sessionKey := range loginResponses {
		//save in session
		userServerUrlByKey[sessionKey.(string)] = url
	}
	return loginResponses, nil
}

type DefaultService struct{}

type DefaultCallParams struct {
	HubKey     string
	ServerArgs [][]interface{}
}

func (h *Auth) AttachToServer(r *http.Request, args *struct{ HubSessionKey, ServerURL, Username, Password string }, reply *struct{ Data string }) error {
	if isHubSessionValid(args.HubSessionKey) {
		serverUsername := args.Username
		serverPass := args.Password

		if autoRelayMode {
			serverUsername = username
			serverPass = pass
		}
		response, _ := executeXMLRPCCall(args.ServerURL, "auth.login", []interface{}{serverUsername, serverPass})
		reply.Data = response.(string)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *DefaultCallParams, reply *struct{ Data map[string]interface{} }) error {
	if isHubSessionValid(args.HubKey) {
		method, _ := NewCodec().NewRequest(r).Method()
		reply.Data = multicastCall(method, args.ServerArgs)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func multicastCall(method string, serverArgs [][]interface{}) map[string]interface{} {
	responses := make(map[string]interface{})
	//Execute the calls concurrently but wait before we get the response from all the servers.
	var wg sync.WaitGroup

	wg.Add(len(serverArgs))

	for _, args := range serverArgs {
		//TODO: check for conversion errors
		url := userServerUrlByKey[args[0].(string)].(string)
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
