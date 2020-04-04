package gateway

import (
	"errors"
	"log"
	"sync"
)

type Multicaster interface {
	Multicast(hubSessionKey, call string, argsByServer map[int64][]interface{}) (*MulticastResponse, error)
}

type multicaster struct {
	client  Client
	session Session
}

func NewMulticaster(client Client, session Session) *multicaster {
	return &multicaster{client, session}
}

func (m *multicaster) Multicast(hubSessionKey, call string, argsByServer map[int64][]interface{}) (*MulticastResponse, error) {
	hubSession := m.session.RetrieveHubSession(hubSessionKey)
	if hubSession == nil {
		log.Printf("HubSession was not found. HubSessionKey: %v", hubSessionKey)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	serverCalls, err := m.generateServerCallRequests(call, hubSession.ServerSessions, argsByServer)
	if err != nil {
		return nil, err
	}
	return executeCallOnServers(serverCalls, m.client), nil
}

type serverCallRequest struct {
	endpoint string
	call     string
	serverID int64
	args     []interface{}
}

type ServerCallResponse struct {
	endpoint string
	call     string
	ServerID int64
	Response interface{}
}

func (m *multicaster) generateServerCallRequests(call string, serverSessions map[int64]*ServerSession, argsByServer map[int64][]interface{}) ([]serverCallRequest, error) {
	result := make([]serverCallRequest, 0, len(argsByServer))

	for serverID, serverArgs := range argsByServer {
		if serverSession, ok := serverSessions[serverID]; ok {
			args := append([]interface{}{serverSession.serverSessionKey}, serverArgs...)
			result = append(result, serverCallRequest{serverSession.serverAPIEndpoint, call, serverID, args})
		} else {
			log.Printf("ServerSession was not found. ServerID: %v", serverID)
			return nil, errors.New("Authentication error: provided session key is invalid")
		}
	}
	return result, nil
}

type MulticastResponse struct {
	SuccessfulResponses, FailedResponses []ServerCallResponse
}

func executeCallOnServers(serverCalls []serverCallRequest, client Client) *MulticastResponse {
	var mutexForSuccesfulResponses = &sync.Mutex{}
	var mutexForFailedResponses = &sync.Mutex{}

	successfulResponses := make([]ServerCallResponse, 0)
	failedResponses := make([]ServerCallResponse, 0)

	var wg sync.WaitGroup
	wg.Add(len(serverCalls))

	for _, serverCall := range serverCalls {
		go func(endpoint, call string, serverID int64, args []interface{}) {
			defer wg.Done()
			response, err := client.ExecuteCall(endpoint, call, args)
			if err != nil {
				log.Printf("Error ocurred in multicast call, serverID: %v, call:%v, error: %v", serverID, call, err)
				mutexForFailedResponses.Lock()
				failedResponses = append(failedResponses, ServerCallResponse{endpoint, call, serverID, err.Error()})
				mutexForFailedResponses.Unlock()
			} else {
				mutexForSuccesfulResponses.Lock()
				successfulResponses = append(successfulResponses, ServerCallResponse{endpoint, call, serverID, response})
				mutexForSuccesfulResponses.Unlock()
			}
		}(serverCall.endpoint, serverCall.call, serverCall.serverID, serverCall.args)
	}
	wg.Wait()
	return &MulticastResponse{successfulResponses, failedResponses}
}
