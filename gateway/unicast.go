package gateway

import (
	"errors"
	"log"
)

type Unicaster interface {
	Unicast(hubSessionKey string, call string, serverID int64, args []interface{}) (interface{}, error)
}

type unicaster struct {
	uyuniServerCallExecutor UyuniServerCallExecutor
	serverSessionRepository ServerSessionRepository
}

func NewUnicaster(uyuniServerCallExecutor UyuniServerCallExecutor, serverSessionRepository ServerSessionRepository) *unicaster {
	return &unicaster{uyuniServerCallExecutor, serverSessionRepository}
}

func (u *unicaster) Unicast(hubSessionKey string, call string, serverID int64, args []interface{}) (interface{}, error) {
	serverSession := u.serverSessionRepository.RetrieveServerSessionByServerID(hubSessionKey, serverID)
	if serverSession == nil {
		log.Printf("ServerSession was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	callArguments := append([]interface{}{serverSession.serverSessionKey}, args...)
	return u.uyuniServerCallExecutor.ExecuteCall(serverSession.serverAPIEndpoint, call, callArguments)
}
