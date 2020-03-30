package service

import (
	"errors"
	"log"
	"sync"
)

type MulticastService struct {
	*service
}

func NewMulticastService(client Client, session Session, hubSumaAPIURL string) *MulticastService {
	return &MulticastService{&service{client: client, session: session, hubSumaAPIURL: hubSumaAPIURL}}
}

func (h *MulticastService) ExecuteMulticastCall(hubSessionKey, path string, serverIDs []int64, serverArgs [][]interface{}) (*MulticastResponse, error) {
	if h.isHubSessionValid(hubSessionKey) {
		serverArgsByURL, err := h.resolveMulticastServerArgs(hubSessionKey, serverIDs, serverArgs)
		if err != nil {
			return nil, err
		}
		return performMulticastCall(path, serverArgsByURL, h.client), nil
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

func (h *MulticastService) resolveMulticastServerArgs(hubSessionKey string, serverIDs []int64, allServerArgs [][]interface{}) ([]multicastServerArgs, error) {
	result := make([]multicastServerArgs, len(serverIDs))
	for i, serverID := range serverIDs {
		serverArgs := make([]interface{}, 0, len(allServerArgs)+1)

		serverSession := h.session.RetrieveServerSessionByServerID(hubSessionKey, serverID)
		if serverSession == nil {
			log.Printf("ServerSessionKey was not found. HubSessionKey: %v, ServerID: %v", hubSessionKey, serverID)
			return nil, errors.New("provided session key is invalid")
		}

		serverArgs = append(serverArgs, serverSession.sessionKey)
		for _, serverArgs := range allServerArgs {
			serverArgs = append(serverArgs, serverArgs[i])
		}
		result[i] = multicastServerArgs{serverSession.url, serverID, serverArgs}
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

func performMulticastCall(method string, args []multicastServerArgs, client Client) *MulticastResponse {
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
