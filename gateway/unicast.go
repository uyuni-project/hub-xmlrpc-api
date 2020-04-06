package gateway

import (
	"errors"
	"log"
)

type Unicaster interface {
	Unicast(request *UnicastRequest) (interface{}, error)
}

type UnicastRequest struct {
	Call          string
	HubSessionKey string
	ServerID      int64
	Args          []interface{}
}

type unicaster struct {
	uyuniServerCallExecutor UyuniServerCallExecutor
	session                 Session
}

func NewUnicaster(uyuniServerCallExecutor UyuniServerCallExecutor, session Session) *unicaster {
	return &unicaster{uyuniServerCallExecutor, session}
}

func (u *unicaster) Unicast(request *UnicastRequest) (interface{}, error) {
	serverSession := u.session.RetrieveServerSessionByServerID(request.HubSessionKey, request.ServerID)
	if serverSession == nil {
		log.Printf("ServerSession was not found. HubSessionKey: %v, ServerID: %v", request.HubSessionKey, request.ServerID)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	callArguments := append([]interface{}{serverSession.serverSessionKey}, request.Args...)
	return u.uyuniServerCallExecutor.ExecuteCall(serverSession.serverAPIEndpoint, request.Call, callArguments)
}
