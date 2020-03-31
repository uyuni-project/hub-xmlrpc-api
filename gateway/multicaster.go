package gateway

import (
	"errors"
	"log"
	"sync"
)

type Multicaster interface {
	Multicast(hubSessionKey, path string, argsByServer map[int64][]interface{}) (*MulticastResponse, error)
}

type MulticastService struct {
	client           Client
	session          Session
	sessionValidator sessionValidator
}

func NewMulticastService(client Client, session Session, sessionValidator sessionValidator) *MulticastService {
	return &MulticastService{client, session, sessionValidator}
}

func (h *MulticastService) Multicast(hubSessionKey, path string, argsByServer map[int64][]interface{}) (*MulticastResponse, error) {
	if h.sessionValidator.isHubSessionValid(hubSessionKey) {
		serverArgsByURL, err := h.resolveMulticastServerArgs(hubSessionKey, argsByServer)
		if err != nil {
			return nil, err
		}
		return executeMulticastCall(path, serverArgsByURL, h.client), nil
	}
	log.Printf("Provided session key is invalid: %v", hubSessionKey)
	//TODO: should we return an error here?
	return nil, nil
}

type multicastServerArgs struct {
	url      string
	serverID int64
	args     []interface{}
}

func (h *MulticastService) resolveMulticastServerArgs(hubSessionKey string, argsByServer map[int64][]interface{}) ([]multicastServerArgs, error) {
	result := make([]multicastServerArgs, 0, len(argsByServer))

	for serverID, serverArgs := range argsByServer {
		serverSession := h.session.RetrieveServerSessionByServerID(hubSessionKey, serverID)
		if serverSession == nil {
			log.Printf("ServerSessionKey was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
			return nil, errors.New("provided session key is invalid")
		}

		args := append([]interface{}{serverSession.sessionKey}, serverArgs...)
		result = append(result, multicastServerArgs{serverSession.url, serverID, args})
	}
	return result, nil
}

type MulticastStateResponse struct {
	ServerIds []int64
	Responses []interface{}
}

type MulticastResponse struct {
	Successful, Failed MulticastStateResponse
}

func executeMulticastCall(method string, args []multicastServerArgs, client Client) *MulticastResponse {
	var mutexForSuccesfulResponses = &sync.Mutex{}
	var mutexForFailedResponses = &sync.Mutex{}

	successfulResponses := make(map[int64]interface{})
	failedResponses := make(map[int64]interface{})

	var wg sync.WaitGroup
	wg.Add(len(args))

	for _, serverArgs := range args {
		go func(url string, serverArgs []interface{}, serverID int64) {
			defer wg.Done()
			response, err := client.ExecuteCall(url, method, serverArgs)
			if err != nil {
				log.Printf("Error ocurred in multicast call, serverID: %v, call:%v, error: %v", serverID, method, err)
				mutexForFailedResponses.Lock()
				failedResponses[serverID] = err.Error()
				mutexForFailedResponses.Unlock()
			} else {
				mutexForSuccesfulResponses.Lock()
				successfulResponses[serverID] = response
				mutexForSuccesfulResponses.Unlock()
			}
		}(serverArgs.url, serverArgs.args, serverArgs.serverID)
	}
	wg.Wait()

	successfulKeys, successfulValues := getServerIDsAndResponses(successfulResponses)
	failedKeys, failedValues := getServerIDsAndResponses(failedResponses)

	return &MulticastResponse{MulticastStateResponse{successfulKeys, successfulValues}, MulticastStateResponse{failedKeys, failedValues}}
}

func getServerIDsAndResponses(in map[int64]interface{}) ([]int64, []interface{}) {
	serverIDs := make([]int64, 0, len(in))
	responses := make([]interface{}, 0, len(in))

	for serverID, response := range in {
		serverIDs = append(serverIDs, serverID)
		responses = append(responses, response)
	}
	return serverIDs, responses
}
