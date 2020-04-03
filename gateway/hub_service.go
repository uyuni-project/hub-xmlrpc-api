package gateway

import "log"

const listSystemsPath = "system.listSystems"

type HubService interface {
	ListServerIDs(hubSessionKey string) ([]int64, error)
}

type HubServiceImpl struct {
	client         Client
	hubAPIEndpoint string
}

func NewHubServiceImpl(client Client, hubAPIEndpoint string) *HubServiceImpl {
	return &HubServiceImpl{client, hubAPIEndpoint}
}

func (h *HubServiceImpl) ListServerIDs(hubSessionKey string) ([]int64, error) {
	systemList, err := h.client.ExecuteCall(h.hubAPIEndpoint, listSystemsPath, []interface{}{hubSessionKey})
	if err != nil {
		log.Printf("Error occured while retrieving the list of serverIDs: %v", err)
		return nil, err
	}
	systemsSlice := systemList.([]interface{})

	systemIDs := make([]int64, len(systemsSlice))
	for i, system := range systemsSlice {
		systemIDs[i] = system.(map[string]interface{})["id"].(int64)
	}
	return systemIDs, nil
}
