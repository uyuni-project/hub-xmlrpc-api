package uyuni

import (
	"errors"
	"log"
	"sync"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

const (
	listSystemsPath                = "system.listSystems"
	listSystemsWithEntitlementPath = "system.listSystemsWithEntitlement"
	listSystemFQDNsPath            = "system.listFqdns"
	listUserSystemsPath            = "system.listUserSystems"
	listPeripheralServersPath      = "sync.hub.listPeripheralServers"
	systemIDField                  = "id"
	peripheralServerEntitlement    = "peripheral_server"
)

type uyuniTopologyInfoRetriever struct {
	uyuniCallExecutor *uyuniCallExecutor
	useSSL            bool
	serverEndpoints   map[int64]string
	mu                sync.RWMutex
}

func NewUyuniTopologyInfoRetriever(uyuniCallExecutor *uyuniCallExecutor, useSSL bool) *uyuniTopologyInfoRetriever {
	return &uyuniTopologyInfoRetriever{
		uyuniCallExecutor: uyuniCallExecutor,
		useSSL:            useSSL,
		serverEndpoints:   make(map[int64]string),
	}
}

func (h *uyuniTopologyInfoRetriever) RetrieveUserServerIDs(endpoint, sessionKey, username string) ([]int64, error) {
	userServers, err := h.uyuniCallExecutor.ExecuteCall(endpoint, listUserSystemsPath, []interface{}{sessionKey, username})
	if err != nil {
		log.Printf("Error ocurred while trying to login into the user systems: %v", err)
		return nil, err
	}

	allServers, err := h.ListServerIDs(endpoint, sessionKey)
	if err != nil {
		return nil, err
	}

	userServersSlice := userServers.([]interface{})

	serverIDs := make([]int64, 0, len(userServersSlice))
	for _, userSystem := range userServersSlice {
		serverID := userSystem.(map[string]interface{})[systemIDField].(int64)
		if contains(allServers, serverID) {
			serverIDs = append(serverIDs, serverID)
		}
	}
	return serverIDs, nil
}

func contains(s []int64, v int64) bool {
	for _, lv := range s {
		if lv == v {
			return true
		}
	}

	return false
}

func (h *uyuniTopologyInfoRetriever) ListServerIDs(endpoint, sessionKey string) ([]int64, error) {
	peripheralServers, err := h.uyuniCallExecutor.ExecuteCall(endpoint, listPeripheralServersPath, []interface{}{sessionKey})
	if err != nil {
		log.Printf("Error occured while calling listPeripheralServersPath: %v", err)
	} else {
		serversSlice := peripheralServers.([]interface{})
		systemIDs := make([]int64, len(serversSlice))
		h.mu.Lock()
		defer h.mu.Unlock()
		for i, s := range serversSlice {
			server := s.(map[string]interface{})
			id := server["id"].(int64)
			fqdn := server["fqdn"].(string)
			rootCA := server["root_ca"].(string)
			systemIDs[i] = id

			serverEndpoint := constructEndpoint(fqdn, h.useSSL)
			h.serverEndpoints[id] = serverEndpoint
			if rootCA != "" {
				err := h.uyuniCallExecutor.client.AddRootCA([]byte(rootCA))
				if err != nil {
					log.Printf("Error ocurred while adding root CA for serverID %v: %v", id, err)
				}
			}
		}
		return systemIDs, nil
	}

	systemList, err := h.uyuniCallExecutor.ExecuteCall(endpoint, listSystemsWithEntitlementPath, []interface{}{sessionKey, peripheralServerEntitlement})
	if err != nil {
		log.Printf("Error occured while retrieving the list of serverIDs: %v", err)
		return nil, err
	}
	systemsSlice := systemList.([]interface{})
	if len(systemsSlice) == 0 {
		// No entitled servers - fallback to full list, for legacy HUB server
		systemList, err = h.uyuniCallExecutor.ExecuteCall(endpoint, listSystemsPath, []interface{}{sessionKey})
		if err != nil {
			log.Printf("Error occured while retrieving the list of serverIDs: %v", err)
			return nil, err
		}
		systemsSlice = systemList.([]interface{})
	}

	systemIDs := make([]int64, len(systemsSlice))
	for i, system := range systemsSlice {
		systemIDs[i] = system.(map[string]interface{})[systemIDField].(int64)
	}
	return systemIDs, nil
}

func (h *uyuniTopologyInfoRetriever) RetrieveServerAPIEndpoints(endpoint, sessionKey string, serverIDs []int64) (*gateway.RetrieveServerAPIEndpointsResponse, error) {
	serverAPIEndpointByServer := make(map[int64]string)
	failedServers := make(map[int64]string)
	for _, serverID := range serverIDs {
		h.mu.RLock()
		serverAPIEndpoint, ok := h.serverEndpoints[serverID]
		h.mu.RUnlock()
		if !ok {
			var err error
			serverAPIEndpoint, err = h.retrieveServerAPIEndpoint(endpoint, sessionKey, serverID)
			if err != nil {
				failedServers[serverID] = err.Error()
				continue
			}
			h.mu.Lock()
			h.serverEndpoints[serverID] = serverAPIEndpoint
			h.mu.Unlock()
		}
		serverAPIEndpointByServer[serverID] = serverAPIEndpoint
	}
	return &gateway.RetrieveServerAPIEndpointsResponse{serverAPIEndpointByServer, failedServers}, nil
}

func (h *uyuniTopologyInfoRetriever) retrieveServerAPIEndpoint(endpoint, sessionKey string, serverID int64) (string, error) {
	//if more than one FQDN is retrieve, we keep the first one and discard the rest
	response, err := h.uyuniCallExecutor.ExecuteCall(endpoint, listSystemFQDNsPath, []interface{}{sessionKey, serverID})
	if err != nil {
		log.Printf("Error ocurred when retrieving the system Fqdns for serverID: %v, error:%v", serverID, err)
		return "", err
	}
	return parseFQDN(response, h.useSSL)
}

func parseFQDN(fqdnResponse interface{}, useSSL bool) (string, error) {
	fqdns, ok := fqdnResponse.([]interface{})
	if !ok {
		log.Printf("Error ocurred when parsing the FQDNs of peripheral servers")
		return "", errors.New("Error ocurred when parsing the FQDNs of peripheral servers")
	}
	if len(fqdns) < 1 {
		log.Printf("Error ocurred when retrieving the FQDNs of peripheral servers: no FQDN found")
		return "", errors.New("Error ocurred when retrieving the FQDNs of peripheral servers: no FQDN found")
	}
	firstFqdn, ok := fqdns[0].(string)
	if !ok {
		log.Printf("Error ocurred when parsing the FQDNs of peripheral servers")
		return "", errors.New("Error ocurred when parsing the FQDNs of peripheral servers")
	}
	return constructEndpoint(firstFqdn, useSSL), nil
}

func constructEndpoint(fqdn string, useSSL bool) string {
	protocol := "http://"
	if useSSL {
		protocol = "https://"
	}
	return protocol + fqdn + "/rpc/api"
}
