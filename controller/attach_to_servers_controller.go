package controller

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type ServerAuthenticationController struct {
	serverAuthenticator gateway.ServerAuthenticator
	responseTransformer multicastResponseTransformer
}

func NewServerAuthenticationController(serverAuthenticator gateway.ServerAuthenticator, responseTransformer multicastResponseTransformer) *ServerAuthenticationController {
	return &ServerAuthenticationController{serverAuthenticator, responseTransformer}
}

type AttachToServersRequest struct {
	HubSessionKey       string
	ServerIDs           []int64
	CredentialsByServer map[int64]*gateway.Credentials
}

func (h *ServerAuthenticationController) AttachToServers(r *http.Request, args *AttachToServersRequest, reply *struct{ Data *MulticastResponse }) error {
	//TODO: what to do with the response?
	attachToServersResponse, err := h.serverAuthenticator.AttachToServers(args.HubSessionKey, args.ServerIDs, args.CredentialsByServer)
	if err != nil {
		log.Printf("Login error: %v", err)
		return err
	}
	reply.Data = h.responseTransformer(attachToServersResponse)
	return nil
}
