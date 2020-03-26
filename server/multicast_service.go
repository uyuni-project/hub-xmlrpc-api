package server

import (
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
		method := removeMulticastNamespace(args.Method)
		serverArgsByURL := h.resolveMulticastServerArgs(args)
		reply.Data = multicastCall(method, serverArgsByURL, h.client)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

type multicastServerArgs struct {
	url      string
	serverID int64
	args     []interface{}
}

func (h *MulticastService) resolveMulticastServerArgs(multicastArgs *MulticastArgs) []multicastServerArgs {
	result := make([]multicastServerArgs, len(multicastArgs.ServerIDs))
	for i, serverID := range multicastArgs.ServerIDs {
		args := make([]interface{}, 0, len(multicastArgs.ServerArgs)+1)

		serverSession := h.session.RetrieveServerSessionByServerID(multicastArgs.HubSessionKey, serverID)
		args = append(args, serverSession.sessionKey)

		for _, serverArgs := range multicastArgs.ServerArgs {
			args = append(args, serverArgs[i])
		}
		result[i] = multicastServerArgs{serverSession.url, serverID, args}
	}
	return result
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

func multicastCall(method string, serverArgs []multicastServerArgs, client Client) MulticastResponse {
	var mutexForSuccesfulResponses = &sync.Mutex{}
	var mutexForFailedResponses = &sync.Mutex{}

	successfulResponses := make(map[int64]interface{})
	failedResponses := make(map[int64]interface{})

	var wg sync.WaitGroup
	wg.Add(len(serverArgs))

	for _, args := range serverArgs {
		go func(url string, args []interface{}, serverId int64) {
			defer wg.Done()
			response, err := client.ExecuteCall(url, method, args)
			if err != nil {
				log.Printf("Call error: %v", err)
				mutexForFailedResponses.Lock()
				failedResponses[serverId] = err.Error()
				mutexForFailedResponses.Unlock()
			} else {
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
