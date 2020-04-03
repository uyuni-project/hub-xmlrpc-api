package gateway

import (
	"errors"
	"log"
)

type Unicaster interface {
	Unicast(hubSessionKey, call string, serverID int64, serverArgs []interface{}) (interface{}, error)
}

type unicaster struct {
	client  Client
	session Session
}

func NewUnicaster(client Client, session Session) *unicaster {
	return &unicaster{client, session}
}

func (h *unicaster) Unicast(hubSessionKey, call string, serverID int64, serverArgs []interface{}) (interface{}, error) {
	serverSession := h.session.RetrieveServerSessionByServerID(hubSessionKey, serverID)
	if serverSession == nil {
		log.Printf("ServerSession was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	callArguments := append([]interface{}{serverSession.serverSessionKey}, serverArgs...)

	return h.client.ExecuteCall(serverSession.serverAPIEndpoint, call, callArguments)
}
