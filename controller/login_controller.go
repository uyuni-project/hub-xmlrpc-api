package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type HubLoginController struct {
	hubLoginer gateway.HubLoginer
}

type LoginRequest struct {
	Username string
	Password string
}

func NewHubLoginController(hubLoginer gateway.HubLoginer) *HubLoginController {
	return &HubLoginController{hubLoginer}
}

func (h *HubLoginController) Login(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.hubLoginer.Login(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubLoginController) LoginWithAuthRelayMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.hubLoginer.LoginWithAuthRelayMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *HubLoginController) LoginWithAutoconnectMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.hubLoginer.LoginWithAutoconnectMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}
