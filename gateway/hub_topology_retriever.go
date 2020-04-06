package gateway

import "log"

type HubTopologyInfoRetriever interface {
	ListServerIDs(hubSessionKey string) ([]int64, error)
}

type hubTopologyInfoRetriever struct {
	uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever
}

func NewHubTopologyInfoRetriever(uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever) *hubTopologyInfoRetriever {
	return &hubTopologyInfoRetriever{uyuniHubTopologyInfoRetriever}
}

func (h *hubTopologyInfoRetriever) ListServerIDs(hubSessionKey string) ([]int64, error) {
	serverIDs, err := h.uyuniHubTopologyInfoRetriever.ListServerIDs(hubSessionKey)
	if err != nil {
		log.Printf("Error occured while retrieving the list of serverIDs: %v", err)
		return nil, err
	}
	return serverIDs, nil
}
