package server

import (
	"log"
	"net/http"
)

type ListArgs struct{ Args []interface{} }

type DefaultService struct {
	client Client
}

func NewDefaultService(client Client) *DefaultService {
	return &DefaultService{client: client}
}

func (h *DefaultService) DefaultMethod(r *http.Request, args *ListArgs, reply *struct{ Data interface{} }) error {
	method, _ := NewCodec().NewRequest(r).Method()
	response, err := h.client.ExecuteCallToHub(method, args.Args)
	if err != nil {
		log.Printf("Call error: %v", err)
	}
	reply.Data = response
	return nil
}
