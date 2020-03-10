package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/session"
)

type Hub struct{}

func (h *Hub) ListServerIds(r *http.Request, args *struct{ HubSessionKey string }, reply *struct{ Data []int64 }) error {
	hubSessionKey := args.HubSessionKey

	if isHubSessionValid(hubSessionKey) {
		systemList, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listSystems", []interface{}{hubSessionKey})
		if err != nil {
			log.Printf("Login error: %v", err)
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
	hubSessionKey, err := h.loginToHub(args.Username, args.Password, session.LOGIN_MANUAL_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *Hub) LoginWithAutoconnectMode(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	hubSessionKey, err := h.loginToHub(args.Username, args.Password, session.LOGIN_AUTOCONNECT_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *Hub) LoginWithAuthRelayMode(r *http.Request, args *struct{ Username, Password string }, reply *struct{ Data string }) error {
	hubSessionKey, err := h.loginToHub(args.Username, args.Password, session.LOGIN_RELAY_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *Hub) loginToHub(username, password string, loginMode int) (string, error) {
	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "auth.login", []interface{}{username, password})
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", errors.New(err.Error())
	}
	hubSessionKey := response.(string)
	apiSession.SetHubSessionKey(hubSessionKey, username, password, loginMode)

	if loginMode == session.LOGIN_AUTOCONNECT_MODE {
		err := loginIntoUserSystems(hubSessionKey, username, password)
		if err != nil {
			log.Printf("Call error: %v", err)
		}
	}
	return hubSessionKey, nil
}

type AttachToServersArgs struct {
	HubSessionKey string
	ServerIDs     []int64
	Usernames     []string
	Passwords     []string
}

func (h *Hub) AttachToServers(r *http.Request, args *AttachToServersArgs, reply *struct{ Data []error }) error {
	if isHubSessionValid(args.HubSessionKey) {
		usernames := args.Usernames
		passwords := args.Passwords

		if apiSession.GetLoginMode(args.HubSessionKey) == session.LOGIN_RELAY_MODE {
			serverUsername, serverPassword := apiSession.GetUsernameAndPassword(args.HubSessionKey)
			usernames = make([]string, len(args.ServerIDs))
			passwords = make([]string, len(args.ServerIDs))

			for i := range args.ServerIDs {
				usernames[i] = serverUsername
				passwords[i] = serverPassword
			}
		}
		loginIntoSystems(args.HubSessionKey, args.ServerIDs, args.Usernames, args.Passwords)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func loginIntoUserSystems(hubSessionKey, username, password string) error {
	userSystems, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listUserSystems", []interface{}{hubSessionKey, username})
	if err != nil {
		log.Printf("Login error: %v", err)
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

func loginIntoSystems(hubSessionKey string, serverIDs []int64, usernames, passwords []string) (MulticastResponse, error) {
	loginIntoSystemsArgs, serverURLByServerID, _ := resolveLoginIntoSystemsArgs(hubSessionKey, serverIDs, usernames, passwords)
	responses := multicastCall("auth.login", loginIntoSystemsArgs)
	successfulResponses := responses.Successfull

	//save in session
	for i, serverID := range successfulResponses.ServerIds {
		apiSession.SetServerSessionInfo(hubSessionKey, serverID, serverURLByServerID[serverID], successfulResponses.Responses[i].(string))
	}
	return responses, nil
}

func resolveLoginIntoSystemsArgs(hubSessionKey string, serverIDs []int64, usernames, passwords []string) ([]MulticastServerArgs, map[int64]string, error) {
	multicastServerArgs := make([]MulticastServerArgs, len(serverIDs))
	serverURLByServerID := make(map[int64]string)

	for i, serverID := range serverIDs {
		url, err := retrieveServerXMLRPCApiURL(hubSessionKey, serverID)
		if err != nil {
			log.Printf("Login error: %v", err)
			//TODO: what to do with failing servers?
		} else {
			serverURLByServerID[serverID] = url
			multicastServerArgs[i] = MulticastServerArgs{url, serverID, []interface{}{usernames[i], passwords[i]}}
		}
	}
	return multicastServerArgs, serverURLByServerID, nil
}

func retrieveServerXMLRPCApiURL(hubSessionKey string, serverID int64) (string, error) {
	//TODO: we should deal with cases when we have more than one fqdn
	response, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "system.listFqdns", []interface{}{hubSessionKey, serverID})
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	//TODO: check for casting errors.
	//TODO: check the fqdn array is not empty
	firstFqdn := response.([]interface{})[0].(string)
	return "http://" + firstFqdn + "/rpc/api", nil
}

func isHubSessionValid(hubSessionKey string) bool {
	isValid, err := executeXMLRPCCall(conf.Hub.SUMA_API_URL, "auth.isSessionKeyValid", []interface{}{hubSessionKey})
	if err != nil {
		log.Printf("Login error: %v", err)
		apiSession.RemoveHubSessionKey(hubSessionKey)
		return false
	}
	return isValid.(bool)
}

func areAllArgumentsOfSameLength(allArrays [][]interface{}) bool {
	if len(allArrays) <= 1 {
		return true
	}
	lengthToCompare := len(allArrays[0])
	for _, array := range allArrays {
		if lengthToCompare != len(array) {
			return false
		}
	}
	return true
}
