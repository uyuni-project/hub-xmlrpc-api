package gateway

import "log"

const listSystemsPath = "system.listSystems"

type HubService interface {
	ListServerIDs(hubSessionKey string) ([]int64, error)
}

type HubServiceImpl struct {
	client           Client
	hubAPIEndpoint   string
	sessionValidator sessionValidator
}

func NewHubServiceImpl(client Client, hubAPIEndpoint string, sessionValidator sessionValidator) *HubServiceImpl {
	return &HubServiceImpl{client, hubAPIEndpoint, sessionValidator}
}

func (h *HubServiceImpl) ListServerIDs(hubSessionKey string) ([]int64, error) {
	if h.sessionValidator.isHubSessionKeyValid(hubSessionKey) {
		systemList, err := h.client.ExecuteCall(h.hubAPIEndpoint, listSystemsPath, []interface{}{hubSessionKey})
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