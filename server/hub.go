package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/chiaradiamarcelo/hub_xmlrpc_api/session"
)

type Hub struct{}

func (h *Hub) ListServerIds(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data []int64 }) error {
	//TODO: parse
	hubSessionKey := args.ArgsList[0].(string)

	if isHubSessionValid(hubSessionKey) {
		systemList, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listSystems", []interface{}{hubSessionKey})
		if err != nil {
			log.Println("Login error: %v", err)
		}
		systemsSlice := systemList.([]interface{})

		systemIDs := make([]int64, len(systemsSlice))
		for i, system := range systemsSlice {
			systemIDs[i] = system.(map[string]interface{})["id"].(int64)
		}
		reply.Data = systemIDs
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func (h *Hub) Login(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data string }) error {
	//TODO: parse
	username := args.ArgsList[0].(string)
	password := args.ArgsList[1].(string)

	hubSessionKey, err := h.loginToHub(username, password, session.LOGIN_MANUAL_MODE)
	if err != nil {
		log.Println("Login error: %v", err)
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *Hub) LoginWithAutoconnectMode(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data string }) error {
	//TODO: parse
	username := args.ArgsList[0].(string)
	password := args.ArgsList[1].(string)

	hubSessionKey, err := h.loginToHub(username, password, session.LOGIN_AUTOCONNECT_MODE)
	if err != nil {
		log.Println("Login error: %v", err)
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *Hub) LoginWithAuthRelayMode(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data string }) error {
	//TODO: parse
	username := args.ArgsList[0].(string)
	password := args.ArgsList[1].(string)

	hubSessionKey, err := h.loginToHub(username, password, session.LOGIN_RELAY_MODE)
	if err != nil {
		log.Println("Login error: %v", err)
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *Hub) loginToHub(username, password string, loginMode int) (string, error) {
	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "auth.login", []interface{}{username, password})
	if err != nil {
		log.Println("Login error: %v", err)
	}
	hubSessionKey := response.(string)
	apiSession.SetHubSessionKey(hubSessionKey, username, password, loginMode)

	if loginMode == session.LOGIN_AUTOCONNECT_MODE || conf.AutoConnectMode {
		err := loginIntoUserSystems(hubSessionKey, username, password)
		if err != nil {
			log.Println("Call error: %v", err)
		}
	}
	return hubSessionKey, nil
}

func (h *Hub) AttachToServers(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data []error }) error {
	//TODO: parse
	hubSessionKey, serverIDs, usernamesArgs, passwordsArgs := parseAttachToServerArguments(args.ArgsList)

	if isHubSessionValid(hubSessionKey) {
		usernames := usernamesArgs
		passwords := passwordsArgs

		if apiSession.GetLoginMode(hubSessionKey) == session.LOGIN_RELAY_MODE || conf.RelayMode {
			serverUsername, serverPassword := apiSession.GetUsernameAndPassword(hubSessionKey)
			usernames = make([]string, len(serverIDs))
			passwords = make([]string, len(serverIDs))

			for i := range serverIDs {
				usernames[i] = serverUsername
				passwords[i] = serverPassword
			}
		}
		loginIntoSystems(hubSessionKey, serverIDs, usernames, passwords)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func parseAttachToServerArguments(argsList []interface{}) (string, []int64, []string, []string) {
	hubSessionKey := argsList[0].(string)

	serverIDs := make([]int64, len(argsList[1].([]interface{})))
	for i, elem := range argsList[1].([]interface{}) {
		serverIDs[i] = elem.(int64)
	}

	var usernames, passwords []string
	if len(argsList) == 4 {
		usernames := make([]string, len(argsList[2].([]interface{})))
		for i, elem := range argsList[2].([]interface{}) {
			usernames[i] = elem.(string)
		}

		passwords := make([]string, len(argsList[3].([]interface{})))
		for i, elem := range argsList[3].([]interface{}) {
			passwords[i] = elem.(string)
		}
	}
	return hubSessionKey, serverIDs, usernames, passwords
}

func loginIntoUserSystems(hubSessionKey, username, password string) error {
	userSystems, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listUserSystems", []interface{}{hubSessionKey, username})
	if err != nil {
		log.Println("Login error: %v", err)
		return err
	}
	userSystemsSlice := userSystems.([]interface{})

	serverIDs := make([]int64, len(userSystemsSlice))
	usernames := make([]string, len(userSystemsSlice))
	passwords := make([]string, len(userSystemsSlice))

	for i, userSystem := range userSystemsSlice {
		serverIDs[i] = userSystem.(map[string]interface{})["id"].(int64)
		usernames[i] = username
		passwords[i] = password
	}

	loginIntoSystems(hubSessionKey, serverIDs, usernames, passwords)
	return nil
}

func loginIntoSystems(hubSessionKey string, serverIDs []int64, usernames, passwords []string) error {
	//TODO: check usernames, passwords and serverIDs have the same size
	var wg sync.WaitGroup
	wg.Add(len(serverIDs))
	for i, serverID := range serverIDs {
		go func(serverID int64, username, password string) {
			defer wg.Done()
			//TODO: we should get the server URL from the 'userSystem'
			url := conf.ServerURLByServerID[strconv.FormatInt(serverID, 10)]
			loginIntoSystem(hubSessionKey, serverID, url, username, password)
		}(serverID, usernames[i], passwords[i])
	}
	wg.Wait()
	return nil
}

func loginIntoSystem(hubSessionKey string, serverID int64, serverURL, username, password string) error {
	response, err := executeXMLRPCCall(serverURL, "auth.login", []interface{}{username, password})
	if err != nil {
		log.Println("Login error: %v", err)
		return err
	}
	//save in session
	apiSession.SetServerSessionInfo(hubSessionKey, serverID, serverURL, response.(string))
	return nil
}
func isHubSessionValid(hubSessionKey string) bool {
	isValid, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "auth.isSessionKeyValid", []interface{}{hubSessionKey})
	if err != nil {
		log.Println("Login error: %v", err)
		apiSession.RemoveHubSessionKey(hubSessionKey)
		return false
	}
	return isValid.(bool)
}
