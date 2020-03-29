package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/uyuni-project/hub-xmlrpc-api/service"
)

type UnicastController struct {
	service *service.UnicastService
}

func NewUnicastController(service *service.UnicastService) *UnicastController {
	return &UnicastController{service}
}

type UnicastRequest struct {
	Method        string
	HubSessionKey string
	ServerID      int64
	ServerArgs    []interface{}
}

func (h *UnicastController) DefaultMethod(r *http.Request, args *UnicastRequest, reply *struct{ Data interface{} }) error {
	method := removeUnicastNamespace(args.Method)
	response, err := h.service.ExecuteUnicastCall(args.HubSessionKey, method, args.ServerID, args.ServerArgs)
	if err != nil {
		log.Printf("Call error: %v", err)
		return err
	}
	reply.Data = response
	return nil
}

func removeUnicastNamespace(method string) string {
	//TODO: removing multicast namespace
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}
