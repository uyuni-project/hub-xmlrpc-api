package server

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

type Multicast struct{}

type MulticastArgs struct {
	HubKey     string
	ServerIds  []int64
	ServerArgs [][]interface{}
}

func (h *Multicast) DefaultMethod(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data []interface{} }) error {
	//TODO: parse
	hubKey, serverIds, serverArgs := parseArgs(args.ArgsList)

	if IsHubSessionValid(hubKey) {
		method, err := NewCodec().NewRequest(r).Method()
		//TODO: HACK for removing multicast namespace
		method = removeMulticastNamespace(method)
		if err != nil {
			log.Println("Call error: %v", err)
		}
		//TODO: check args.ServerArgs lists have the same size
		serverArgsByURL := make(map[string][]interface{})

		for i, serverID := range serverIds {
			out := make([]interface{}, len(serverArgs)+1)

			for j, serverArgs := range serverArgs {
				out[j+1] = serverArgs[i]
			}
			url, sessionKey := apiSession.GetServerSessionInfoByServerID(hubKey, serverID)
			out[0] = sessionKey
			serverArgsByURL[url] = out
		}
		reply.Data = multicastCall(method, serverArgsByURL)
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func parseArgs(argsList []interface{}) (string, []int64, [][]interface{}) {
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

func removeMulticastNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}

func multicastCall(method string, serverArgsByURL map[string][]interface{}) []interface{} {
	responses := make([]interface{}, len(serverArgsByURL))

	var wg sync.WaitGroup
	wg.Add(len(serverArgsByURL))

	i := 0
	for url, args := range serverArgsByURL {
		go func(url string, args []interface{}, i int) {
			defer wg.Done()
			response, err := executeXMLRPCCall(url, method, args)
			if err != nil {
				log.Println("Call error: %v", err)
			}
			responses[i] = response
			log.Printf("Response: %s\n", response)
		}(url, args, i)
		i++
	}
	wg.Wait()
	return responses
}
