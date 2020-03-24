package server

import (
	"log"
	"net/http"
	"strings"
)

type UnicastService struct {
	client  Client
	session Session
}

func NewUnicastService(client Client, session Session) *UnicastService {
	return &UnicastService{client: client, session: session}
}

func (h *UnicastService) DefaultMethod(r *http.Request, args *UnicastArgs, reply *struct{ Data interface{} }) error {
	if h.session.IsHubSessionValid(args.HubSessionKey) {
		method := removeUnicastNamespace(args.Method)
		argumentsForCall := make([]interface{}, 0, len(args.ServerArgs)+1)
		url, sessionKey := h.session.GetServerSessionInfoByServerID(args.HubSessionKey, args.ServerID)
		argumentsForCall = append(argumentsForCall, sessionKey)
		argumentsForCall = append(argumentsForCall, args.ServerArgs...)

		response, err := h.client.ExecuteCall(url, method, argumentsForCall)
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
