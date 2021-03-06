package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type HubTopologyController struct {
	hubService gateway.TopologyInfoRetriever
}

func (h *HubTopologyController) ListServerIDs(r *http.Request, args *struct{ HubSessionKey string }, reply *struct{ Data []int64 }) error {
	serverIDs, err := h.hubService.ListServerIDs(args.HubSessionKey)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = serverIDs
	return nil
}

func NewHubTopologyController(topologyInfoRetriever gateway.TopologyInfoRetriever) *HubTopologyController {
	return &HubTopologyController{topologyInfoRetriever}
}
