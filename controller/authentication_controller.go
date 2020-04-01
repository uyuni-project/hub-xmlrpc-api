package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type AuthenticationController struct {
	authenticator gateway.Authenticator
}

type LoginRequest struct {
	Username string
	Password string
}

func NewAuthenticationController(authenticator gateway.Authenticator) *AuthenticationController {
	return &AuthenticationController{authenticator}
}

func (h *AuthenticationController) Login(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.authenticator.Login(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *AuthenticationController) LoginWithAutoconnectMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.authenticator.LoginWithAutoconnectMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *AuthenticationController) LoginWithAuthRelayMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.authenticator.LoginWithAuthRelayMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *AuthenticationController) AttachToServers(r *http.Request, args *MulticastRequest, reply *struct{ Data []error }) error {
	//TODO: what to do with the response?
	_, err := h.authenticator.AttachToServers(args.HubSessionKey, args.ArgsByServer)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	return nil
}
