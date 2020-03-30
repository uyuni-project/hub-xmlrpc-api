package service

import (
	"errors"
	"log"
)

type UnicastService struct {
	*service
}

func NewUnicastService(client Client, session Session, hubSumaAPIURL string) *UnicastService {
	return &UnicastService{&service{client: client, session: session, hubSumaAPIURL: hubSumaAPIURL}}
}

func (h *UnicastService) ExecuteUnicastCall(hubSessionKey, path string, serverID int64, serverArgs []interface{}) (interface{}, error) {
	if h.isHubSessionValid(hubSessionKey) {
		serverSession := h.session.RetrieveServerSessionByServerID(hubSessionKey, serverID)
		if serverSession == nil {
			log.Printf("ServerSessionKey was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
			return nil, errors.New("provided session key is invalid")
		}

		argumentsForCall := append([]interface{}{serverSession.sessionKey}, serverArgs...)

		return h.client.ExecuteCall(serverSession.url, path, argumentsForCall)
	}
	log.Printf("Provided session key is invalid: %v", hubSessionKey)
	//TODO: should we return an error here?
	return nil, nil
}
