package server

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type UnicastService struct {
	*service
}

func NewUnicastService(client Client, session Session, hubSumaAPIURL string) *UnicastService {
	return &UnicastService{&service{client: client, session: session, hubSumaAPIURL: hubSumaAPIURL}}
}

type UnicastArgs struct {
	Method        string
	HubSessionKey string
	ServerID      int64
	ServerArgs    []interface{}
}

func (h *UnicastService) DefaultMethod(r *http.Request, args *UnicastArgs, reply *struct{ Data interface{} }) error {
	if h.isHubSessionValid(args.HubSessionKey) {
		serverSession := h.session.RetrieveServerSessionByServerID(args.HubSessionKey, args.ServerID)
		if serverSession == nil {
			log.Printf("ServerSessionKey was not found. HubSessionKey: %v, ServerID: %v", args.HubSessionKey, args.ServerID)
			return errors.New("provided session key is invalid")
		}

		argumentsForCall := make([]interface{}, 0, len(args.ServerArgs)+1)
		argumentsForCall = append(argumentsForCall, serverSession.sessionKey)
		argumentsForCall = append(argumentsForCall, args.ServerArgs...)

		method := removeUnicastNamespace(args.Method)

		response, err := h.client.ExecuteCall(serverSession.url, method, argumentsForCall)
		if err != nil {
			log.Printf("Call error: %v", err)
			return err
		}
		reply.Data = response
	} else {
		log.Printf("Provided session key is invalid: %v", args.HubSessionKey)
		//TODO: should we return an error here?
	}
	return nil
}

func removeUnicastNamespace(method string) string {
	//TODO: removing multicast namespace
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}
