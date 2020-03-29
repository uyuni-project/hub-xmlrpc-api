package controller

import (
	"net/http"
	"strings"

	"github.com/uyuni-project/hub-xmlrpc-api/service"
)

type MulticastController struct {
	service *service.MulticastService
}

func NewMulticastController(service *service.MulticastService) *MulticastController {
	return &MulticastController{service}
}

type MulticastRequest struct {
	Method        string
	HubSessionKey string
	ServerIDs     []int64
	ServerArgs    [][]interface{}
}

type MulticastResponse struct {
	Successfull, Failed MulticastStateResponse
}

type MulticastStateResponse struct {
	ServerIds []int64
	Responses []interface{}
}

func (h *MulticastController) DefaultMethod(r *http.Request, args *MulticastRequest, reply *struct{ Data *service.MulticastResponse }) error {
	method := removeMulticastNamespace(args.Method)
	response, err := h.service.ExecuteMulticastCall(args.HubSessionKey, method, args.ServerIDs, args.ServerArgs)
	if err != nil {
		return err
	}
	reply.Data = response
	return nil
}

func removeMulticastNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}
