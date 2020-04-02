package gateway

import (
	"errors"
	"log"
	"sync"
)

type Multicaster interface {
	Multicast(hubSessionKey, call string, argsByServer map[int64][]interface{}) (*MulticastResponse, error)
}

type MulticastService struct {
	client           Client
	session          Session
	sessionValidator sessionValidator
}

func NewMulticastService(client Client, session Session, sessionValidator sessionValidator) *MulticastService {
	return &MulticastService{client, session, sessionValidator}
}

func (m *MulticastService) Multicast(hubSessionKey, call string, argsByServer map[int64][]interface{}) (*MulticastResponse, error) {
	if m.sessionValidator.isHubSessionKeyValid(hubSessionKey) {
		serverCalls, err := m.appendServerSessionKeyToServerArgs(hubSessionKey, argsByServer)
		if err != nil {
			return nil, err
		}
		return executeCallOnServers(call, serverCalls, m.client), nil
	}
	log.Printf("Provided session key is invalid: %v", hubSessionKey)
	//TODO: what error message return?
	return nil, errors.New("Authentication error: provided session key is invalid")
}

type serverCall struct {
	serverID       int64
	serverEndpoint string
	serverArgs     []interface{}
}

func (m *MulticastService) appendServerSessionKeyToServerArgs(hubSessionKey string, argsByServer map[int64][]interface{}) ([]serverCall, error) {
	result := make([]serverCall, 0, len(argsByServer))

	serverSessions := m.session.RetrieveServerSessions(hubSessionKey)

	for serverID, serverArgs := range argsByServer {
		if serverSession, ok := serverSessions[serverID]; ok {
			args := append([]interface{}{serverSession.serverSessionKey}, serverArgs...)
			result = append(result, serverCall{serverID, serverSession.serverAPIEndpoint, args})
		} else {
			log.Printf("ServerSession was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
			return nil, errors.New("Authentication error: provided session key is invalid")
		}
	}
	return result, nil
}

type MulticastResponse struct {
	SuccessfulResponses, FailedResponses map[int64]interface{}
}

func executeCallOnServers(call string, serverCalls []serverCall, client Client) *MulticastResponse {
	var mutexForSuccesfulResponses = &sync.Mutex{}
	var mutexForFailedResponses = &sync.Mutex{}

	successfulResponses := make(map[int64]interface{})
	failedResponses := make(map[int64]interface{})

	var wg sync.WaitGroup
	wg.Add(len(serverCalls))

	for _, serverCall := range serverCalls {
		go func(serverEndpoint string, serverID int64, serverArgs []interface{}) {
			defer wg.Done()
			response, err := client.ExecuteCall(serverEndpoint, call, serverArgs)
			if err != nil {
				log.Printf("Error ocurred in multicast call, serverID: %v, call:%v, error: %v", serverID, call, err)
				mutexForFailedResponses.Lock()
				failedResponses[serverID] = err.Error()
				mutexForFailedResponses.Unlock()
			} else {
				mutexForSuccesfulResponses.Lock()
				successfulResponses[serverID] = response
				mutexForSuccesfulResponses.Unlock()
			}
		}(serverCall.serverEndpoint, serverCall.serverID, serverCall.serverArgs)
	}
	wg.Wait()
	return &MulticastResponse{successfulResponses, failedResponses}
}
