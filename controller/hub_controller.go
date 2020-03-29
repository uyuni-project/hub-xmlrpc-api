package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/service"
)

type HubController struct {
	service *service.HubService
}

func NewHubController(hubService *service.HubService) *HubController {
	return &HubController{hubService}
}

type LoginRequest struct {
	Username string
	Password string
}

func (h *HubController) Login(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.service.Login(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubController) LoginWithAutoconnectMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.service.LoginWithAutoconnectMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubController) LoginWithAuthRelayMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.service.LoginWithAuthRelayMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubController) AttachToServers(r *http.Request, args *MulticastRequest, reply *struct{ Data []error }) error {
	_, err := h.service.AttachToServers(args.HubSessionKey, args.ServerIDs, args.ServerArgs)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	return nil
}

func (h *HubController) ListServerIds(r *http.Request, args *struct{ HubSessionKey string }, reply *struct{ Data []int64 }) error {
	serverIDs, err := h.service.ListServerIds(args.HubSessionKey)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = serverIDs
	return nil
}
