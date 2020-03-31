package controller

import (
	"net/http"
	"strings"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type MulticastController struct {
	multicaster gateway.Multicaster
}

type MulticastRequest struct {
	Method        string
	HubSessionKey string
	ArgsByServer  map[int64][]interface{}
}

func NewMulticastController(multicaster gateway.Multicaster) *MulticastController {
	return &MulticastController{multicaster}
}

func (h *MulticastController) Multicast(r *http.Request, args *MulticastRequest, reply *struct{ Data *gateway.MulticastResponse }) error {
	method := removeMulticastNamespace(args.Method)
	response, err := h.multicaster.Multicast(args.HubSessionKey, method, args.ArgsByServer)
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
