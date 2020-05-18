package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type HubLoginController struct {
	hubLoginer          gateway.HubLoginer
	responseTransformer multicastResponseTransformer
}

type LoginRequest struct {
	Username string
	Password string
}

func NewHubLoginController(hubLoginer gateway.HubLoginer, responseTransformer multicastResponseTransformer) *HubLoginController {
	return &HubLoginController{hubLoginer, responseTransformer}
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

type LoginWithAutoconnectModeResponse struct {
	SessionKey         string
	Successful, Failed MulticastStateResponse
}

func (h *HubLoginController) LoginWithAutoconnectMode(r *http.Request, args *LoginRequest, reply *struct {
	Data *LoginWithAutoconnectModeResponse
}) error {
	loginResponse, err := h.hubLoginer.LoginWithAutoconnectMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	attachToServersResponse := h.responseTransformer(loginResponse.AttachToServersResponse)
	reply.Data = &LoginWithAutoconnectModeResponse{loginResponse.HubSessionKey, attachToServersResponse.Successful, attachToServersResponse.Failed}
	return nil
}
