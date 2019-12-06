package server

import (
	"log"
	"net/http"
	"strings"
)

type Unicast struct{}

func (h *Unicast) DefaultMethod(r *http.Request, args *struct{ ArgsList []interface{} }, reply *struct{ Data interface{} }) error {
	//TODO: HACK parse
	hubKey, serverID, serverArgs := parseUnicastArgs(args.ArgsList)

	if IsHubSessionValid(hubKey) {
		method, err := NewCodec().NewRequest(r).Method()
		//TODO: HACK for removing multicast namespace
		method = removeUnicastNamespace(method)
		if err != nil {
			log.Println("Call error: %v", err)
		}
		argumentsForCall := make([]interface{}, len(serverArgs)+1)
		url, sessionKey := apiSession.GetServerSessionInfoByServerID(hubKey, serverID)
		argumentsForCall[0] = sessionKey

		response, err := executeXMLRPCCall(url, method, argumentsForCall)
		if err != nil {
			log.Println("Call error: %v", err)
		}
		reply.Data = response
	} else {
		log.Println("Hub session invalid error")
	}
	return nil
}

func parseUnicastArgs(argsList []interface{}) (string, int64, []interface{}) {
	//TODO: HACK
	hubKey := argsList[0].(string)
	serverID := argsList[1].(int64)
	serverArgs := argsList[2:len(argsList)]
	return hubKey, serverID, serverArgs
}

func removeUnicastNamespace(method string) string {
	//TODO: HACK for removing multicast namespace
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}
