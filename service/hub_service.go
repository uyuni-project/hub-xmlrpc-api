package service

import (
	"errors"
	"log"
)

const (
	LOGIN_PATH             = "auth.login"
	LIST_SYSTEMS_PATH      = "system.listSystems"
	LIST_USER_SYSTEMS_PATH = "system.listUserSystems"
	LIST_SYSTEM_FQDNS_PATH = "system.listFqdns"

	LOGIN_MANUAL_MODE      = iota // 0
	LOGIN_RELAY_MODE              // 1
	LOGIN_AUTOCONNECT_MODE        // 2
)

type HubService struct {
	*service
}

func NewHubService(client Client, session Session, hubSumaAPIURL string) *HubService {
	return &HubService{&service{client: client, session: session, hubSumaAPIURL: hubSumaAPIURL}}
}

func (h *HubService) Login(username, password string) (string, error) {
	hubSessionKey, err := h.loginToHub(username, password, LOGIN_MANUAL_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return hubSessionKey, nil
}

func (h *HubService) LoginWithAutoconnectMode(username, password string) (string, error) {
	hubSessionKey, err := h.loginToHub(username, password, LOGIN_AUTOCONNECT_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return hubSessionKey, nil
}

func (h *HubService) LoginWithAuthRelayMode(username, password string) (string, error) {
	hubSessionKey, err := h.loginToHub(username, password, LOGIN_RELAY_MODE)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return hubSessionKey, nil
}

func (h *HubService) AttachToServers(hubSessionKey string, serverIDs []int64, serverArgs [][]interface{}) (*MulticastResponse, error) {
	if h.isHubSessionValid(hubSessionKey) {
		usernames := make([]interface{}, len(serverIDs))
		passwords := make([]interface{}, len(serverIDs))

		hubSession := h.session.RetrieveHubSession(hubSessionKey)
		if hubSession == nil {
			log.Printf("HubSessionKey was not found: %v", hubSessionKey)
			//TODO: what error should we return here?
			return nil, errors.New("provided session key is invalid")
		}

		if hubSession.loginMode == LOGIN_RELAY_MODE {
			for i := range serverIDs {
				usernames[i] = hubSession.username
				passwords[i] = hubSession.password
			}
		} else {
			usernames = serverArgs[0]
			passwords = serverArgs[1]
		}
		return h.loginIntoSystems(hubSessionKey, serverIDs, usernames, passwords)
	}
	log.Printf("Provided session key is invalid: %v", hubSessionKey)
	//TODO: should we return an error here?
	return nil, nil
}

func (h *HubService) ListServerIds(hubSessionKey string) ([]int64, error) {
	if h.isHubSessionValid(hubSessionKey) {
		systemList, err := h.client.ExecuteCall(h.hubSumaAPIURL, LIST_SYSTEMS_PATH, []interface{}{hubSessionKey})
		if err != nil {
			log.Printf("Login error: %v", err)
			return nil, err
		}
		systemsSlice := systemList.([]interface{})

		systemIDs := make([]int64, len(systemsSlice))
		for i, system := range systemsSlice {
			systemIDs[i] = system.(map[string]interface{})["id"].(int64)
		}
		return systemIDs, nil
	}
	log.Printf("Provided session key is invalid: %v", hubSessionKey)
	//TODO: should we return an error here?
	return nil, nil
}

func (h *HubService) loginToHub(username, password string, loginMode int) (string, error) {
	response, err := h.client.ExecuteCall(h.hubSumaAPIURL, LOGIN_PATH, []interface{}{username, password})
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	hubSessionKey := response.(string)
	h.session.SaveHubSession(hubSessionKey, &HubSession{username, password, loginMode})

	if loginMode == LOGIN_AUTOCONNECT_MODE {
		err := h.loginIntoUserSystems(hubSessionKey, username, password)
		if err != nil {
			//TODO: should we return an error? retry the login or what?
			log.Printf("Call error: %v", err)
		}
	}
	return hubSessionKey, nil
}

func (h *HubService) loginIntoUserSystems(hubSessionKey, username, password string) error {
	userSystems, err := h.client.ExecuteCall(h.hubSumaAPIURL, LIST_USER_SYSTEMS_PATH, []interface{}{hubSessionKey, username})
	if err != nil {
		log.Printf("Error ocurred while trying to logiin into the user systems: %v", err)
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

	//TODO: what to do with the response here?
	_, err = h.loginIntoSystems(hubSessionKey, serverIDs, usernames, passwords)
	return err
}

func (h *HubService) loginIntoSystems(hubSessionKey string, serverIDs []int64, usernames, passwords []interface{}) (*MulticastResponse, error) {
	//TODO: what to do with the error here?
	loginIntoSystemsArgs, serverURLByServerID, _ := h.resolveLoginIntoSystemsArgs(hubSessionKey, serverIDs, usernames, passwords)
	responses := performMulticastCall(LOGIN_PATH, loginIntoSystemsArgs, h.client)
	successfulResponses := responses.Successful

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
		url, err := h.retrieveServerAPIURL(hubSessionKey, serverID)
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

func (h *HubService) retrieveServerAPIURL(hubSessionKey string, serverID int64) (string, error) {
	//TODO: we should deal with cases when we have more than one fqdn
	response, err := h.client.ExecuteCall(h.hubSumaAPIURL, LIST_SYSTEM_FQDNS_PATH, []interface{}{hubSessionKey, serverID})
	if err != nil {
		log.Printf("Error ocurred when retrieving the system Fqdns for serverID: %v, error:%v", serverID, err)
		return "", err
	}
	firstFqdn := response.([]interface{})[0].(string)
	if firstFqdn == "mch-server-slave-1.tf.local" {
		return "http://192.168.122.188/rpc/api", nil
	}
	return "http://192.168.122.196/rpc/api", nil
}
