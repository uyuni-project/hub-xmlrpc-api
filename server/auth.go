package server

import (
	"log"
	"net/http"
	"strconv"
)

type Auth struct{}

func (h *Auth) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data map[string]interface{} }) error {
	sessionkeys := make(map[string]interface{})

	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "auth.login", []interface{}{args.Username, args.Password})
	if err != nil {
		log.Println("Login error: %v", err)
	}
	//TODO: remove when Abids PR is merged
	hubSessionKey := response.(string)
	apiSession.SetHubSessionKey(hubSessionKey)

	sessionkeys["hubSessionKey"] = hubSessionKey

	if conf.RelayMode {
		//save credentials in session
		apiSession.SetUsernameAndPassword(args.Username, args.Password)

		if conf.AutoConnectMode {
			serverSessionKeys, err := loginIntoUserSystems(hubSessionKey, args.Username, args.Password)
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
		url := conf.ServerURLByServerID[strconv.FormatInt(systemID, 10)]
		serverArgsByURL[url] = []interface{}{username, password}
	}
	//TODO: reuse loginIntoSystem method
	loginResponses := multicastCall("auth.login", serverArgsByURL)

	//save in session
	for url, sessionKey := range loginResponses {
		apiSession.AddServerURLforServerKey(url, sessionKey.(string))
	}
	return loginResponses, nil
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
			serverUsername, serverPass = apiSession.GetUsernameAndPassword()
		}
		reply.Data, _ = loginIntoSystem(args.ServerID, serverUsername, serverPass)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func loginIntoSystem(serverID int64, username, password string) (string, error) {
	serverURL := getServerURLFromServerID(serverID)
	response, err := executeXMLRPCCall(serverURL, "auth.login", []interface{}{username, password})

	if err != nil {
		return "", err
	}
	//save in session
	apiSession.AddServerURLforServerKey(serverURL, response.(string))
	return response.(string), nil
}
