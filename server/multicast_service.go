package server

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

type MulticastService struct {
	client  Client
	session Session
}

func NewMulticastService(client Client, session Session) *MulticastService {
	return &MulticastService{client: client, session: session}
}

type MulticastArgs struct {
	HubSessionKey string
	ServerIDs     []int64
	ServerArgs    [][]interface{}
}

func (h *MulticastService) DefaultMethod(r *http.Request, args *MulticastArgs, reply *struct{ Data MulticastResponse }) error {
	if !areAllArgumentsOfSameLength(args.ServerArgs) {
		return FaultInvalidParams
	}
	if h.session.IsHubSessionValid(args.HubSessionKey) {
		method, err := NewCodec().NewRequest(r).Method()
		//TODO: removing multicast namespace. We should reuse the same codec we use for the server
		method = removeMulticastNamespace(method)
		if err != nil {
			log.Printf("Call error: %v", err)
		}
		serverArgsByURL := h.resolveMulticastServerArgs(args)
		reply.Data = multicastCall(method, serverArgsByURL, h.client)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

type MulticastServerArgs struct {
	url      string
	serverID int64
	args     []interface{}
}

func (h *MulticastService) resolveMulticastServerArgs(multicastArgs *MulticastArgs) []MulticastServerArgs {
	multicastServerArgs := make([]MulticastServerArgs, len(multicastArgs.ServerIDs))
	for i, serverID := range multicastArgs.ServerIDs {
		args := make([]interface{}, 0, len(multicastArgs.ServerArgs)+1)

		url, sessionKey := h.session.GetServerSessionInfoByServerID(multicastArgs.HubSessionKey, serverID)
		args = append(args, sessionKey)

		for _, serverArgs := range multicastArgs.ServerArgs {
			args = append(args, serverArgs[i])
		}
		multicastServerArgs[i] = MulticastServerArgs{url, serverID, args}
	}
	return multicastServerArgs
}

func removeMulticastNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}

type MulticastResponse struct {
	Successfull, Failed MulticastStateResponse
}

type MulticastStateResponse struct {
	Responses []interface{}
	ServerIds []int64
}

func multicastCall(method string, serverArgs []MulticastServerArgs, client Client) MulticastResponse {
	var mutexForSuccesfulResponses = &sync.Mutex{}
	var mutexForFailedResponses = &sync.Mutex{}

	successfulResponses := make(map[int64]interface{})
	failedResponses := make(map[int64]interface{})

	var wg sync.WaitGroup
	wg.Add(len(serverArgs))

	for _, args := range serverArgs {
		go func(url string, args []interface{}, serverId int64) {
			defer wg.Done()
			response, err := client.ExecuteCallWithURL(url, method, args)
			if err != nil {
				log.Printf("Call error: %v", err)
				mutexForFailedResponses.Lock()
				failedResponses[serverId] = err.Error()
				mutexForFailedResponses.Unlock()
			} else {
				log.Printf("Response: %s\n", response)
				mutexForSuccesfulResponses.Lock()
				successfulResponses[serverId] = response
				mutexForSuccesfulResponses.Unlock()
			}
		}(args.url, args.args, args.serverID)
	}
	wg.Wait()

	successfulKeys, successfulValues := getKeysAndValuesFromMap(successfulResponses)
	failedKeys, failedValues := getKeysAndValuesFromMap(failedResponses)

	return MulticastResponse{MulticastStateResponse{successfulValues, successfulKeys}, MulticastStateResponse{failedValues, failedKeys}}
}

func getKeysAndValuesFromMap(in map[int64]interface{}) ([]int64, []interface{}) {
	keys := make([]int64, 0, len(in))
	values := make([]interface{}, 0, len(in))

	for key, value := range in {
		keys = append(keys, key)
		values = append(values, value)
	}
	return keys, values
}

func areAllArgumentsOfSameLength(allArrays [][]interface{}) bool {
	if len(allArrays) <= 1 {
		return true
	}
	lengthToCompare := len(allArrays[0])
	for _, array := range allArrays {
		if lengthToCompare != len(array) {
			return false
		}
	}
	return true
}
