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
	argsByServer, err := resolveArgsByServer(args.HubSessionKey, args.ServerIDs, args.ServerArgs)
	response, err := h.service.ExecuteMulticastCall(args.HubSessionKey, method, argsByServer)
	if err != nil {
		return err
	}
	reply.Data = response
	return nil
}

func resolveArgsByServer(hubSessionKey string, serverIDs []int64, allServerArgs [][]interface{}) (map[int64][]interface{}, error) {
	result := make(map[int64][]interface{})
	for i, serverID := range serverIDs {
		args := make([]interface{}, 0, len(allServerArgs)+1)

		for _, serverArgs := range allServerArgs {
			args = append(args, serverArgs[i])
		}
		result[serverID] = args
	}
	return result, nil
}

func removeMulticastNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}
