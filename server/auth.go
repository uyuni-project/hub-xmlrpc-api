package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"
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
		serverURL := getServerURLFromServerID(args.ServerID)
		reply.Data, _ = loginIntoSystem(args.ServerID, serverURL, serverUsername, serverPass)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func loginIntoUserSystems(hubSessionKey, username, password string) ([]map[string]interface{}, error) {
	userSystems, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listUserSystems", []interface{}{hubSessionKey, username})
	if err != nil {
		return nil, err
	}
	userSystemArr := userSystems.([]interface{})
	out := make([]map[string]interface{}, len(userSystemArr))

	var wg sync.WaitGroup
	wg.Add(len(userSystemArr))
	for i, userSystem := range userSystemArr {
		go func(i int, userSystem interface{}) {
			defer wg.Done()
			//TODO: we should get the server URL from the 'userSystem'
			serverID := userSystem.(map[string]interface{})["id"].(int64)
			url := conf.ServerURLByServerID[strconv.FormatInt(serverID, 10)]
			sessionKey, _ := loginIntoSystem(serverID, url, username, password)
			out[i] = map[string]interface{}{"url": url, "sessionKey": sessionKey, "serverID": serverID}
		}(i, userSystem)
	}
	wg.Wait()
	return out, nil
}

func loginIntoSystem(serverID int64, serverURL, username, password string) (string, error) {
	response, err := executeXMLRPCCall(serverURL, "auth.login", []interface{}{username, password})
	if err != nil {
		return "", err
	}
	//save in session
	apiSession.AddServerURLforServerKey(serverURL, response.(string))
	return response.(string), nil
}

func getServerURLFromServerID(serverID int64) string {
	//TODO:
	return conf.ServerURLByServerID[strconv.FormatInt(serverID, 10)]
}
