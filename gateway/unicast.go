package gateway

import (
	"errors"
	"log"
)

type Unicaster interface {
	Unicast(hubSessionKey string, call string, serverID int64, args []interface{}) (interface{}, error)
}

type unicaster struct {
	uyuniCallExecutor       UyuniCallExecutor
	serverSessionRepository ServerSessionRepository
}

func NewUnicaster(uyuniCallExecutor UyuniCallExecutor, serverSessionRepository ServerSessionRepository) *unicaster {
	return &unicaster{uyuniCallExecutor, serverSessionRepository}
}

func (u *unicaster) Unicast(hubSessionKey string, call string, serverID int64, args []interface{}) (interface{}, error) {
	serverSession := u.serverSessionRepository.RetrieveServerSessionByServerID(hubSessionKey, serverID)
	if serverSession == nil {
		log.Printf("ServerSession was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	callArguments := append([]interface{}{serverSession.serverSessionKey}, args...)
	return u.uyuniCallExecutor.ExecuteCall(serverSession.serverAPIEndpoint, call, callArguments)
}
