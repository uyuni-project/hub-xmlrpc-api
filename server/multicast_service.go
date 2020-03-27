package server

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
)

type MulticastService struct {
	*service
}

func NewMulticastService(client Client, session Session, hubSumaAPIURL string) *MulticastService {
	return &MulticastService{&service{client: client, session: session, hubSumaAPIURL: hubSumaAPIURL}}
}

type MulticastArgs struct {
	Method        string
	HubSessionKey string
	ServerIDs     []int64
	ServerArgs    [][]interface{}
}

func (h *MulticastService) DefaultMethod(r *http.Request, args *MulticastArgs, reply *struct{ Data MulticastResponse }) error {
	if h.isHubSessionValid(args.HubSessionKey) {
		serverArgsByURL, err := h.resolveMulticastServerArgs(args)
		if err != nil {
			return err
		}

		method := removeMulticastNamespace(args.Method)
		reply.Data = multicastCall(method, serverArgsByURL, h.client)
	} else {
		log.Printf("Provided session key is invalid: %v", args.HubSessionKey)
		//TODO: should we return an error here?
	}
	return nil
}

type multicastServerArgs struct {
	url      string
	serverID int64
	args     []interface{}
}

func (h *MulticastService) resolveMulticastServerArgs(args *MulticastArgs) ([]multicastServerArgs, error) {
	result := make([]multicastServerArgs, len(args.ServerIDs))
	for i, serverID := range args.ServerIDs {
		serverArgs := make([]interface{}, 0, len(args.ServerArgs)+1)

		serverSession := h.session.RetrieveServerSessionByServerID(args.HubSessionKey, serverID)
		if serverSession == nil {
			log.Printf("ServerSessionKey was not found. HubSessionKey: %v, ServerID: %v", args.HubSessionKey, serverID)
			return nil, errors.New("provided session key is invalid")
		}

		serverArgs = append(serverArgs, serverSession.sessionKey)
		for _, serverArgs := range args.ServerArgs {
			serverArgs = append(serverArgs, serverArgs[i])
		}
		result[i] = multicastServerArgs{serverSession.url, serverID, serverArgs}
	}
	return result, nil
}

func removeMulticastNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}

type MulticastStateResponse struct {
	ServerIds []int64
	Responses []interface{}
}

type MulticastResponse struct {
	Successfull, Failed MulticastStateResponse
}

func multicastCall(method string, args []multicastServerArgs, client Client) MulticastResponse {
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

	return MulticastResponse{MulticastStateResponse{successfulKeys, successfulValues}, MulticastStateResponse{failedKeys, failedValues}}
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
