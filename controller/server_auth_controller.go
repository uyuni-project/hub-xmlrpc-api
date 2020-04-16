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

type AttachToServersRequest struct {
	HubSessionKey       string
	ServerIDs           []int64
	CredentialsByServer map[int64]*gateway.Credentials
}

func (h *ServerAuthenticationController) AttachToServers(r *http.Request, args *AttachToServersRequest, reply *struct{ Data []error }) error {
	//TODO: what to do with the response?
	_, err := h.serverAuthenticator.AttachToServers(args.HubSessionKey, args.ServerIDs, args.CredentialsByServer)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	return nil
}
