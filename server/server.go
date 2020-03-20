package server

import (
	"log"
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
)

type ListArgs struct{ Args []interface{} }

type DefaultService struct {
	client *client.Client
}

func NewDefaultService(client *client.Client) *DefaultService {
	return &DefaultService{client: client}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *ListArgs, reply *struct{ Data interface{} }) error {
	method, _ := NewCodec().NewRequest(r).Method()
	response, err := h.client.ExecuteXMLRPCCallToHub(method, args.Args)
	if err != nil {
		log.Printf("Call error: %v", err)
	}
	reply.Data = response
	return nil
}
