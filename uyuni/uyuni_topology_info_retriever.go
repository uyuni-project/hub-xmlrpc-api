package uyuni

import "log"

const (
	listSystemsPath     = "system.listSystems"
	listSystemFQDNsPath = "system.listFqdns"
	listUserSystemsPath = "system.listUserSystems"
	systemIDField       = "id"
)

type uyuniTopologyInfoRetriever struct {
	uyuniCallExecutor *uyuniCallExecutor
}

func NewUyuniTopologyInfoRetriever(uyuniCallExecutor *uyuniCallExecutor) *uyuniTopologyInfoRetriever {
	return &uyuniTopologyInfoRetriever{uyuniCallExecutor}
}

func (h *uyuniTopologyInfoRetriever) RetrieveUserServerIDs(endpoint, sessionKey, username string) ([]int64, error) {
	userServers, err := h.uyuniCallExecutor.ExecuteCall(endpoint, listUserSystemsPath, []interface{}{sessionKey, username})
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

func (h *uyuniTopologyInfoRetriever) ListServerIDs(endpoint, sessionKey string) ([]int64, error) {
	systemList, err := h.uyuniCallExecutor.ExecuteCall(endpoint, listSystemsPath, []interface{}{sessionKey})
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

func (h *uyuniTopologyInfoRetriever) RetrieveServerAPIEndpoints(endpoint, sessionKey string, serverIDs []int64) (map[int64]string, error) {
	serverAPIEndpointByServer := make(map[int64]string)
	for _, serverID := range serverIDs {
		serverAPIEndpoint, err := h.retrieveServerAPIEndpoint(endpoint, sessionKey, serverID)
		if err != nil {
			//TODO: what to do with failing servers?
		} else {
			serverAPIEndpointByServer[serverID] = serverAPIEndpoint
		}
	}
	return serverAPIEndpointByServer, nil
}

func (h *uyuniTopologyInfoRetriever) retrieveServerAPIEndpoint(endpoint, sessionKey string, serverID int64) (string, error) {
	//TODO: we should deal with cases when we have more than one fqdn
	response, err := h.uyuniCallExecutor.ExecuteCall(endpoint, listSystemFQDNsPath, []interface{}{sessionKey, serverID})
	if err != nil {
		log.Printf("Error ocurred when retrieving the system Fqdns for serverID: %v, error:%v", serverID, err)
		return "", err
	}
	//TODO: check for casting errors.
	//TODO: check the fqdn array is not empty
	firstFqdn := response.([]interface{})[0].(string)
	return "http://" + firstFqdn + "/rpc/api", nil
}
