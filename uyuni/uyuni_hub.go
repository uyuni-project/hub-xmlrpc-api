package uyuni

import (
	"log"
)

//Hub authenticator
type uyuniHubAuthenticator struct {
	uyuniServerAuthenticator *uyuniServerAuthenticator
	hubAPIEndpoint           string
}

func NewUyuniHubAuthenticator(uyuniServerAuthenticator *uyuniServerAuthenticator, hubAPIEndpoint string) *uyuniHubAuthenticator {
	return &uyuniHubAuthenticator{uyuniServerAuthenticator, hubAPIEndpoint}
}

func (h *uyuniHubAuthenticator) Login(username, password string) (string, error) {
	response, err := h.uyuniServerAuthenticator.Login(h.hubAPIEndpoint, username, password)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (h *uyuniHubAuthenticator) Logout(sessionKey string) error {
	err := h.uyuniServerAuthenticator.Logout(h.hubAPIEndpoint, sessionKey)
	if err != nil {
		return err
	}
	return nil
}

//Hub topology information retriever
const (
	listSystemsPath     = "system.listSystems"
	listSystemFQDNsPath = "system.listFqdns"
	listUserSystemsPath = "system.listUserSystems"
	systemIDField       = "id"
)

type uyuniHubTopologyInfoRetriever struct {
	uyuniHubCallExecutor *uyuniHubCallExecutor
}

func NewUyuniHubTopologyInfoRetriever(uyuniHubCallExecutor *uyuniHubCallExecutor) *uyuniHubTopologyInfoRetriever {
	return &uyuniHubTopologyInfoRetriever{uyuniHubCallExecutor}
}

func (h *uyuniHubTopologyInfoRetriever) RetrieveUserServerIDs(sessionKey, username string) ([]int64, error) {
	userServers, err := h.uyuniHubCallExecutor.ExecuteCall(listUserSystemsPath, []interface{}{sessionKey, username})
	if err != nil {
		log.Printf("Error ocurred while trying to login into the user systems: %v", err)
		return nil, err
	}
	userServersSlice := userServers.([]interface{})

	serverIDs := make([]int64, 0, len(userServersSlice))
	for _, userSystem := range userServersSlice {
		serverID := userSystem.(map[string]interface{})[systemIDField].(int64)
		serverIDs = append(serverIDs, serverID)
	}
	return serverIDs, nil
}

func (h *uyuniHubTopologyInfoRetriever) ListServerIDs(sessionKey string) ([]int64, error) {
	systemList, err := h.uyuniHubCallExecutor.ExecuteCall(listSystemsPath, []interface{}{sessionKey})
	if err != nil {
		log.Printf("Error occured while retrieving the list of serverIDs: %v", err)
		return nil, err
	}
	systemsSlice := systemList.([]interface{})

	systemIDs := make([]int64, len(systemsSlice))
	for i, system := range systemsSlice {
		systemIDs[i] = system.(map[string]interface{})[systemIDField].(int64)
	}
	return systemIDs, nil
}

func (h *uyuniHubTopologyInfoRetriever) RetrieveServerAPIEndpoints(sessionKey string, serverIDs []int64) (map[int64]string, error) {
	serverAPIEndpointByServer := make(map[int64]string)
	for _, serverID := range serverIDs {
		serverAPIEndpoint, err := h.retrieveServerAPIEndpoint(sessionKey, serverID)
		if err != nil {
			//TODO: what to do with failing servers?
		} else {
			serverAPIEndpointByServer[serverID] = serverAPIEndpoint
		}
	}
	return serverAPIEndpointByServer, nil
}

func (h *uyuniHubTopologyInfoRetriever) retrieveServerAPIEndpoint(sessionKey string, serverID int64) (string, error) {
	//TODO: we should deal with cases when we have more than one fqdn
	response, err := h.uyuniHubCallExecutor.ExecuteCall(listSystemFQDNsPath, []interface{}{sessionKey, serverID})
	if err != nil {
		log.Printf("Error ocurred when retrieving the system Fqdns for serverID: %v, error:%v", serverID, err)
		return "", err
	}
	//TODO: check for casting errors.
	//TODO: check the fqdn array is not empty
	firstFqdn := response.([]interface{})[0].(string)
	return "http://" + firstFqdn + "/rpc/api", nil
}

//Call executor
type uyuniHubCallExecutor struct {
	uyuniServerCallExecutor *uyuniServerCallExecutor
	hubAPIEndpoint          string
}

func NewUyuniHubCallExecutor(uyuniServerCallExecutor *uyuniServerCallExecutor, hubAPIEndpoint string) *uyuniHubCallExecutor {
	return &uyuniHubCallExecutor{uyuniServerCallExecutor, hubAPIEndpoint}
}

func (h *uyuniHubCallExecutor) ExecuteCall(call string, args []interface{}) (interface{}, error) {
	response, err := h.uyuniServerCallExecutor.ExecuteCall(h.hubAPIEndpoint, call, args)
	if err != nil {
		return "", err
	}
	return response, nil
}
