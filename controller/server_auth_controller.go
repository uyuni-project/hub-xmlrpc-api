package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type ServerAuthenticationController struct {
	serverAuthenticator gateway.ServerAuthenticator
}

func NewServerAuthenticationController(serverAuthenticator gateway.ServerAuthenticator) *ServerAuthenticationController {
	return &ServerAuthenticationController{serverAuthenticator}
}

func (h *ServerAuthenticationController) AttachToServers(r *http.Request, args *gateway.AttachToServersRequest, reply *struct{ Data []error }) error {
	//TODO: what to do with the response?
	_, err := h.serverAuthenticator.AttachToServers(args)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	return nil
}
