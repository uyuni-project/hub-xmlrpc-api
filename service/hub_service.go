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

func (h *HubService) AttachToServers(hubSessionKey string, argsByServerID map[int64][]interface{}) (*MulticastResponse, error) {
	if h.isHubSessionValid(hubSessionKey) {
		hubSession := h.session.RetrieveHubSession(hubSessionKey)
		if hubSession == nil {
			log.Printf("HubSessionKey was not found: %v", hubSessionKey)
			//TODO: what error should we return here?
			return nil, errors.New("provided session key is invalid")
		}

		credentialsByServerID := argsByServerID
		if hubSession.loginMode == LOGIN_RELAY_MODE {
			credentialsByServerID = make(map[int64][]interface{})
			for serverID := range argsByServerID {
				credentialsByServerID[serverID] = []interface{}{hubSession.username, hubSession.password}
			}
		}
		return h.loginIntoSystems(hubSessionKey, credentialsByServerID)
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

	credentialsByServerID := make(map[int64][]interface{})
	for _, userSystem := range userSystemsSlice {
		serverID := userSystem.(map[string]interface{})["id"].(int64)
		credentialsByServerID[serverID] = []interface{}{username, password}
	}

	//TODO: what to do with the response here?
	_, err = h.loginIntoSystems(hubSessionKey, credentialsByServerID)
	return err
}

func (h *HubService) loginIntoSystems(hubSessionKey string, credentialsByServerID map[int64][]interface{}) (*MulticastResponse, error) {
	//TODO: what to do with the error here?
	loginIntoSystemsArgs, serverURLByServerID, _ := h.resolveLoginIntoSystemsArgs(hubSessionKey, credentialsByServerID)
	responses := performMulticastCall(LOGIN_PATH, loginIntoSystemsArgs, h.client)
	successfulResponses := responses.Successful

	//save in session
	for i, serverID := range successfulResponses.ServerIds {
		h.session.SaveServerSession(hubSessionKey, serverID, &ServerSession{serverURLByServerID[serverID], successfulResponses.Responses[i].(string)})
	}
	return responses, nil
}

func (h *HubService) resolveLoginIntoSystemsArgs(hubSessionKey string, credentialsByServerID map[int64][]interface{}) ([]multicastServerArgs, map[int64]string, error) {
	multicastArgs := make([]multicastServerArgs, 0, len(credentialsByServerID))
	serverURLByServerID := make(map[int64]string)

	for serverID, credentials := range credentialsByServerID {
		url, err := h.retrieveServerAPIURL(hubSessionKey, serverID)
		if err != nil {
			log.Printf("Login error: %v", err)
			//TODO: what to do with failing servers?
		} else {
			serverURLByServerID[serverID] = url
			multicastArgs = append(multicastArgs, multicastServerArgs{url, serverID, credentials})
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
	//TODO: check for casting errors.
	//TODO: check the fqdn array is not empty
	firstFqdn := response.([]interface{})[0].(string)
	return "http://" + firstFqdn + "/rpc/api", nil
}
