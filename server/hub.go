package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Hub struct{}

func (h *Hub) ListServerIds(r *http.Request, args *struct{ HubSessionKey string }, reply *struct{ Data []int64 }) error {
	if apiSession.IsHubSessionValid(args.HubSessionKey) {
		systemList, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listSystems", []interface{}{args.HubSessionKey})
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

func (h *Hub) Login(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "auth.login", []interface{}{args.Username, args.Password})
	if err != nil {
		log.Println("Login error: %v", err)
	}
	hubSessionKey := response.(string)
	apiSession.SetHubSessionKey(hubSessionKey, args.Username, args.Password)

	if conf.RelayMode && conf.AutoConnectMode {
		err := loginIntoUserSystems(hubSessionKey, args.Username, args.Password)
		if err != nil {
			log.Println("Call error: %v", err)
		}
	}
	reply.Data = hubSessionKey
	return nil
}

type AttachToServerArgs struct {
	HubSessionKey string
	ServerIDs     []int64
	Usernames     []string
	Passwords     []string
}

func (h *Hub) AttachToServers(r *http.Request, args *AttachToServerArgs, reply *struct{ Data []error }) error {
	if apiSession.IsHubSessionValid(args.HubSessionKey) {
		usernames := args.Usernames
		passwords := args.Passwords

		if conf.RelayMode {
			serverUsername, serverPassword := apiSession.GetUsernameAndPassword(args.HubSessionKey)
			usernames = make([]string, len(args.ServerIDs))
			passwords = make([]string, len(args.ServerIDs))

			for i := range args.ServerIDs {
				usernames[i] = serverUsername
				passwords[i] = serverPassword
			}

		}
		loginIntoSystems(args.HubSessionKey, args.ServerIDs, usernames, passwords)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
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
