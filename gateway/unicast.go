package gateway

import (
	"errors"
	"log"
)

type Unicaster interface {
	Unicast(hubSessionKey, call string, serverID int64, serverArgs []interface{}) (interface{}, error)
}

type UnicastService struct {
	client           Client
	session          Session
	sessionValidator sessionValidator
}

func NewUnicastService(client Client, session Session, sessionValidator sessionValidator) *UnicastService {
	return &UnicastService{client, session, sessionValidator}
}

func (h *UnicastService) Unicast(hubSessionKey, call string, serverID int64, serverArgs []interface{}) (interface{}, error) {
	if h.sessionValidator.isHubSessionKeyValid(hubSessionKey) {
		serverSession := h.session.RetrieveServerSessionByServerID(hubSessionKey, serverID)
		if serverSession == nil {
			log.Printf("ServerSession was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
			return nil, errors.New("provided session key is invalid")
		}

		callArguments := append([]interface{}{serverSession.serverSessionKey}, serverArgs...)

		return h.client.ExecuteCall(serverSession.serverAPIEndpoint, call, callArguments)
	}
	log.Printf("Provided session key is invalid: %v", hubSessionKey)
	//TODO: should we return an error here?
	return nil, nil
}
