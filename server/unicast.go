package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/uyuni-project/hub-xmlrpc-api/session"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

type UnicastArgs struct {
	HubSessionKey string
	ServerID      int64
	ServerArgs    []interface{}
}

type Unicast struct {
	client     *client.Client
	apiSession *session.ApiSession
}

func NewUnicastService(client *client.Client, apiSession *session.ApiSession) *Unicast {
	return &Unicast{client: client, apiSession: apiSession}
}

func (h *Unicast) DefaultMethod(r *http.Request, args *UnicastArgs, reply *struct{ Data interface{} }) error {
	if h.apiSession.IsHubSessionValid(args.HubSessionKey, h.client) {
		method, err := NewCodec().NewRequest(r).Method()
		//TODO: removing multicast namespace. We should reuse the same codec we use for the server
		method = removeUnicastNamespace(method)
		if err != nil {
			log.Printf("Call error: %v", err)
		}
		argumentsForCall := make([]interface{}, 0, len(args.ServerArgs)+1)
		url, sessionKey := h.apiSession.GetServerSessionInfoByServerID(args.HubSessionKey, args.ServerID)
		argumentsForCall = append(argumentsForCall, sessionKey)
		argumentsForCall = append(argumentsForCall, args.ServerArgs...)

		response, err := h.client.ExecuteXMLRPCCallWithURL(url, method, argumentsForCall)
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
