package server

import (
	"errors"
	"log"
	"net/http"
)

type HubService struct {
	*service
}

func NewHubService(client Client, session Session, hubSumaAPIURL string) *HubService {
	return &HubService{&service{client: client, session: session, hubSumaAPIURL: hubSumaAPIURL}}
}

type LoginArgs struct {
	Username string
	Password string
}

func (h *HubService) Login(r *http.Request, args *LoginArgs, reply *struct{ Data string }) error {
	hubSessionKey, err := h.loginToHub(args.Username, args.Password, LOGIN_MANUAL_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubService) LoginWithAutoconnectMode(r *http.Request, args *LoginArgs, reply *struct{ Data string }) error {
	hubSessionKey, err := h.loginToHub(args.Username, args.Password, LOGIN_AUTOCONNECT_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubService) LoginWithAuthRelayMode(r *http.Request, args *LoginArgs, reply *struct{ Data string }) error {
	hubSessionKey, err := h.loginToHub(args.Username, args.Password, LOGIN_RELAY_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubService) AttachToServers(r *http.Request, args *MulticastArgs, reply *struct{ Data []error }) error {
	if h.isHubSessionValid(args.HubSessionKey) {
		usernames := make([]interface{}, len(args.ServerIDs))
		passwords := make([]interface{}, len(args.ServerIDs))

		hubSession := h.session.RetrieveHubSession(args.HubSessionKey)
		if hubSession.loginMode == LOGIN_RELAY_MODE {
			for i := range args.ServerIDs {
				usernames[i] = hubSession.username
				passwords[i] = hubSession.password
			}
		} else {
			usernames = args.ServerArgs[0]
			passwords = args.ServerArgs[1]
		}
		h.loginIntoSystems(args.HubSessionKey, args.ServerIDs, usernames, passwords)
	} else {
		log.Println("Provided session key is invalid.")
		return errors.New("provided session key is invalid")
	}
	return nil
}

func (h *HubService) ListServerIds(r *http.Request, args *struct{ HubSessionKey string }, reply *struct{ Data []int64 }) error {
	if h.isHubSessionValid(args.HubSessionKey) {
		systemList, err := h.client.ExecuteCall(h.hubSumaAPIURL, "system.listSystems", []interface{}{args.HubSessionKey})
		if err != nil {
			log.Printf("Login error: %v", err)
			return err
		}
		systemsSlice := systemList.([]interface{})

		systemIDs := make([]int64, len(systemsSlice))
		for i, system := range systemsSlice {
			systemIDs[i] = system.(map[string]interface{})["id"].(int64)
		}
		reply.Data = systemIDs
	} else {
		log.Println("Provided session key is invalid.")
		return errors.New("Provided session key is invalid")
	}
	return nil
}

func (h *HubService) loginToHub(username, password string, loginMode int) (string, error) {
	response, err := h.client.ExecuteCall(h.hubSumaAPIURL, "auth.login", []interface{}{username, password})
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", errors.New(err.Error())
	}
	hubSessionKey := response.(string)
	h.session.SaveHubSession(hubSessionKey, &HubSession{username, password, loginMode})

	if loginMode == LOGIN_AUTOCONNECT_MODE {
		err := h.loginIntoUserSystems(hubSessionKey, username, password)
		if err != nil {
			log.Printf("Call error: %v", err)
		}
	}
	return hubSessionKey, nil
}

func (h *HubService) loginIntoUserSystems(hubSessionKey, username, password string) error {
	userSystems, err := h.client.ExecuteCall(h.hubSumaAPIURL, "system.listUserSystems", []interface{}{hubSessionKey, username})
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	userSystemsSlice := userSystems.([]interface{})

	serverIDs := make([]int64, len(userSystemsSlice))
	usernames := make([]interface{}, len(userSystemsSlice))
	passwords := make([]interface{}, len(userSystemsSlice))

	for i, userSystem := range userSystemsSlice {
		serverIDs[i] = userSystem.(map[string]interface{})["id"].(int64)
		usernames[i] = username
		passwords[i] = password
	}

	h.loginIntoSystems(hubSessionKey, serverIDs, usernames, passwords)
	return nil
}

func (h *HubService) loginIntoSystems(hubSessionKey string, serverIDs []int64, usernames, passwords []interface{}) (MulticastResponse, error) {
	loginIntoSystemsArgs, serverURLByServerID, _ := h.resolveLoginIntoSystemsArgs(hubSessionKey, serverIDs, usernames, passwords)
	responses := multicastCall("auth.login", loginIntoSystemsArgs, h.client)
	successfulResponses := responses.Successfull

	//save in session
	for i, serverID := range successfulResponses.ServerIds {
		h.session.SaveServerSession(hubSessionKey, serverID, &ServerSession{serverURLByServerID[serverID], successfulResponses.Responses[i].(string)})
	}
	return responses, nil
}

func (h *HubService) resolveLoginIntoSystemsArgs(hubSessionKey string, serverIDs []int64, usernames, passwords []interface{}) ([]multicastServerArgs, map[int64]string, error) {
	multicastArgs := make([]multicastServerArgs, len(serverIDs))
	serverURLByServerID := make(map[int64]string)

	for i, serverID := range serverIDs {
		url, err := h.retrieveServerXMLRPCApiURL(hubSessionKey, serverID)
		if err != nil {
			log.Printf("Login error: %v", err)
			//TODO: what to do with failing servers?
		} else {
			serverURLByServerID[serverID] = url
			multicastArgs[i] = multicastServerArgs{url, serverID, []interface{}{usernames[i], passwords[i]}}
		}
	}
	return multicastArgs, serverURLByServerID, nil
}

func (h *HubService) retrieveServerXMLRPCApiURL(hubSessionKey string, serverID int64) (string, error) {
	//TODO: we should deal with cases when we have more than one fqdn
	response, err := h.client.ExecuteCall(h.hubSumaAPIURL, "system.listFqdns", []interface{}{hubSessionKey, serverID})
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	//TODO: check for casting errors.
	//TODO: check the fqdn array is not empty
	firstFqdn := response.([]interface{})[0].(string)
	return "http://" + firstFqdn + "/rpc/api", nil
}
