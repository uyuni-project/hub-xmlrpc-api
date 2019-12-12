package server

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

type MulticastService struct{}

type MulticastArgs struct {
	HubSessionKey string
	ServerIDs     []int64
	ServerArgs    [][]interface{}
}

func (h *MulticastService) DefaultMethod(r *http.Request, args *MulticastArgs, reply *struct{ Data MulticastResponse }) error {
	if !areAllArgumentsOfSameLength(args.ServerArgs) {
		return FaultInvalidParams
	}
	if isHubSessionValid(args.HubSessionKey) {
		method, err := NewCodec().NewRequest(r).Method()
		//TODO: removing multicast namespace. We should reuse the same codec we use for the server
		method = removeMulticastNamespace(method)
		if err != nil {
			log.Println("Call error: %v", err)
		}
		serverArgsByURL := resolveMulticastServerArgs(args.HubSessionKey, args.ServerIDs, args.ServerArgs)
		reply.Data = multicastCall(method, serverArgsByURL)
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

func resolveMulticastServerArgs(hubSessionKey string, serverIDs []int64, serversArgs [][]interface{}) []MulticastServerArgs {
	multicastServerArgs := make([]MulticastServerArgs, len(serverIDs))
	for i, serverID := range serverIDs {
		args := make([]interface{}, 0, len(serversArgs)+1)

		url, sessionKey := apiSession.GetServerSessionInfoByServerID(hubSessionKey, serverID)
		args = append(args, sessionKey)

		for _, serverArgs := range serversArgs {
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

func multicastCall(method string, serverArgs []MulticastServerArgs) MulticastResponse {
	var mutexForSuccesfulResponses = &sync.Mutex{}
	var mutexForFailedResponses = &sync.Mutex{}

	successfulResponses := make(map[int64]interface{})
	failedResponses := make(map[int64]interface{})

	var wg sync.WaitGroup
	wg.Add(len(serverArgs))

	for _, args := range serverArgs {
		go func(url string, args []interface{}, serverId int64) {
			defer wg.Done()
			response, err := executeXMLRPCCall(url, method, args)
			if err != nil {
				log.Println("Call error: %v", err)
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
