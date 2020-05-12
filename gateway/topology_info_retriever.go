package gateway

import "log"

type TopologyInfoRetriever interface {
	ListServerIDs(hubSessionKey string) ([]int64, error)
}

type topologyInfoRetriever struct {
	hubAPIEndpoint             string
	uyuniTopologyInfoRetriever UyuniTopologyInfoRetriever
}

func NewTopologyInfoRetriever(hubAPIEndpoint string, uyuniTopologyInfoRetriever UyuniTopologyInfoRetriever) *topologyInfoRetriever {
	return &topologyInfoRetriever{hubAPIEndpoint, uyuniTopologyInfoRetriever}
}

func (h *topologyInfoRetriever) ListServerIDs(hubSessionKey string) ([]int64, error) {
	serverIDs, err := h.uyuniTopologyInfoRetriever.ListServerIDs(h.hubAPIEndpoint, hubSessionKey)
	if err != nil {
		log.Printf("Error occured while retrieving the list of serverIDs: %v", err)
		return nil, err
	}
	return serverIDs, nil
}
