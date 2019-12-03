package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/chiaradiamarcelo/hub_xmlrpc_api/client"
	"github.com/chiaradiamarcelo/hub_xmlrpc_api/config"
	"github.com/gorilla/rpc"
)

var conf = config.New()

type Auth struct{}

//TODO: remove when Abid's PR is merged
var hubSessionKey = ""

//TODO: session
var username = ""
var pass = ""
var userServerUrlByKey = make(map[string]string)

//TODO:WE SHOULD GET THIS FROM SUMA API (ie, on listUserSystems)
var serverUrlByServerId = map[int64]string{1000010000: "http://192.168.122.203/rpc/api"}

func isHubSessionValid(in string) bool {
	//TODO: we should check this on session or through the SUMA api
	return in == hubSessionKey
}

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data map[string]interface{} }) error {
	sessionkeys := make(map[string]interface{})

	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "auth.login", []interface{}{args.Username, args.Password})
	if err != nil {
		log.Println("Login error: %v", err)
	}
	//TODO: remove when Abids PR is merged
	hubSessionKey = response.(string)

	sessionkeys["hubSessionKey"] = hubSessionKey

	if conf.RelayMode {
		//save credentials in session
		username = args.Username
		pass = args.Password

		if conf.AutoConnectMode {
			serverSessionKeys, err := loginIntoUserSystems(hubSessionKey, username, pass)
			if err != nil {
				log.Println("Call error: %v", err)
			}
			sessionkeys["serverSessionKeys"] = serverSessionKeys
		}
	}
	reply.Data = sessionkeys
	return nil
}

func loginIntoUserSystems(hubSessionKey, username, password string) (map[string]interface{}, error) {
	userSystems, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listUserSystems", []interface{}{hubSessionKey, username})
	if err != nil {
		return nil, err
	}
	userSystemArr := userSystems.([]interface{})
	serverArgsByURL := make(map[string][]interface{})

	for _, userSystem := range userSystemArr {
		//TODO: we should get the server URL from the 'userSystem'
		systemID := userSystem.(map[string]interface{})["id"].(int64)
		url := serverUrlByServerId[systemID]
		serverArgsByURL[url] = []interface{}{username, password}
	}
	loginResponses := multicastCall("auth.login", serverArgsByURL)

	//save in session
	for url, sessionKey := range loginResponses {
		userServerUrlByKey[sessionKey.(string)] = url
	}
	return loginResponses, nil
}

func loginIntoSystem(serverID int64, username, password string) (string, error) {
	serverURL := getServerURLFromServerId(serverID)
	response, err := executeXMLRPCCall(serverURL, "auth.login", []interface{}{username, password})

	if err != nil {
		return "", err
	}
	//save in session
	userServerUrlByKey[response.(string)] = serverURL
	return response.(string), nil
}

type AttachToServerArgs struct {
	HubSessionKey      string
	ServerID           int64
	Username, Password string
}

func (h *Auth) AttachToServer(r *http.Request, args *AttachToServerArgs, reply *struct{ Data string }) error {
	if isHubSessionValid(args.HubSessionKey) {
		serverUsername := args.Username
		serverPass := args.Password

		if conf.RelayMode {
			serverUsername = username
			serverPass = pass
		}
		reply.Data, _ = loginIntoSystem(args.ServerID, serverUsername, serverPass)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func getServerURLFromServerId(serverID int64) string {
	//TODO:
	return serverUrlByServerId[serverID]
}

type DefaultService struct{}

type DefaultCallArgs struct {
	HubKey     string
	ServerArgs [][]interface{}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *DefaultCallArgs, reply *struct{ Data map[string]interface{} }) error {
	if isHubSessionValid(args.HubKey) {
		method, _ := NewCodec().NewRequest(r).Method()

		serverArgsByURL := make(map[string][]interface{})

		for _, args := range args.ServerArgs {
			//TODO: support methods that don't need sessionkey
			url := userServerUrlByKey[args[0].(string)]
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
