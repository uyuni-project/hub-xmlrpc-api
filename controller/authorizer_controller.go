package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type AuthorizerController struct {
	authorizer gateway.Authorizer
}

type LoginRequest struct {
	Username string
	Password string
}

func NewAuthorizerController(authorizer gateway.Authorizer) *AuthorizerController {
	return &AuthorizerController{authorizer}
}

func (h *AuthorizerController) Login(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.authorizer.Login(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *AuthorizerController) LoginWithAutoconnectMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.authorizer.LoginWithAutoconnectMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *AuthorizerController) LoginWithAuthRelayMode(r *http.Request, args *LoginRequest, reply *struct{ Data string }) error {
	hubSessionKey, err := h.authorizer.LoginWithAuthRelayMode(args.Username, args.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = hubSessionKey
	return nil
}

func (h *AuthorizerController) AttachToServers(r *http.Request, args *MulticastRequest, reply *struct{ Data []error }) error {
	//TODO: what to do with the response?
	_, err := h.authorizer.AttachToServers(args.HubSessionKey, args.ArgsByServer)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	return nil
}
