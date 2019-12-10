package server

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

type MulticastService struct{}

func (h *MulticastService) DefaultMethod(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data MulticastResponse }) error {
	//TODO: parse
	hubSessionKey, serverIds, serverArgs := parseMulticastArgs(args.ArgsList)

	if isHubSessionValid(hubSessionKey) {
		method, err := NewCodec().NewRequest(r).Method()
		//TODO: removing multicast namespace. We should reuse the same codec we use for the server
		method = removeMulticastNamespace(method)
		if err != nil {
			log.Println("Call error: %v", err)
		}
		//TODO: check args.ServerArgs lists have the same size
		serverArgsByURL := resolveMulticastServerArgs(hubSessionKey, serverIds, serverArgs)
		reply.Data = multicastCall(method, serverArgsByURL)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func parseMulticastArgs(argsList []interface{}) (string, []int64, [][]interface{}) {
	//TODO:
	hubKey := argsList[0].(string)
	serverIDs := make([]int64, len(argsList[1].([]interface{})))
	for i, elem := range argsList[1].([]interface{}) {
		serverIDs[i] = elem.(int64)
	}

	rest := argsList[2:len(argsList)]
	serverArgs := make([][]interface{}, len(rest))

	for i, list := range rest {
		serverArgs[i] = list.([]interface{})
	}
	return hubKey, serverIDs, serverArgs
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
				failedResponses[serverId] = err
			} else {
				log.Printf("Response: %s\n", response)
				successfulResponses[serverId] = response
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
