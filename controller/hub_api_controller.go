package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type HubController struct {
	hubService gateway.HubService
}

func (h *HubController) ListServerIDs(r *http.Request, args *struct{ HubSessionKey string }, reply *struct{ Data []int64 }) error {
	serverIDs, err := h.hubService.ListServerIDs(args.HubSessionKey)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = serverIDs
	return nil
}

func NewHubController(hubService gateway.HubService) *HubController {
	return &HubController{hubService}
}
