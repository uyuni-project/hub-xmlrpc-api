package server

import (
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
		method := removeUnicastNamespace(args.Method)
		serverSession := h.session.RetrieveServerSessionByServerID(args.HubSessionKey, args.ServerID)

		argumentsForCall := make([]interface{}, 0, len(args.ServerArgs)+1)
		argumentsForCall = append(argumentsForCall, serverSession.sessionKey)
		argumentsForCall = append(argumentsForCall, args.ServerArgs...)

		response, err := h.client.ExecuteCall(serverSession.url, method, argumentsForCall)
		if err != nil {
			log.Printf("Call error: %v", err)
			return err
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
