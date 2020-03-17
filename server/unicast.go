package server

import (
	"log"
	"net/http"
	"strings"
)

type UnicastArgs struct {
	HubSessionKey string
	ServerID      int64
	ServerArgs    []interface{}
}

type Unicast struct{}

func (h *Unicast) DefaultMethod(r *http.Request, args *UnicastArgs, reply *struct{ Data interface{} }) error {
	if isHubSessionValid(args.HubSessionKey) {
		method, err := NewCodec().NewRequest(r).Method()
		//TODO: removing multicast namespace. We should reuse the same codec we use for the server
		method = removeUnicastNamespace(method)
		if err != nil {
			log.Printf("Call error: %v", err)
		}
		argumentsForCall := make([]interface{}, len(args.ServerArgs)+1)
		url, sessionKey := apiSession.GetServerSessionInfoByServerID(args.HubSessionKey, args.ServerID)
		argumentsForCall[0] = sessionKey

		response, err := executeXMLRPCCall(url, method, argumentsForCall)
		if err != nil {
			log.Printf("Call error: %v", err)
		}
		reply.Data = response
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func removeUnicastNamespace(method string) string {
	//TODO: removing multicast namespace
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}
