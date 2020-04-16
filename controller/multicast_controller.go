package controller

import (
	"net/http"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type MulticastController struct {
	multicaster         gateway.Multicaster
	responseTransformer multicastResponseTransformer
}
type multicastResponseTransformer func(multicastResponse *gateway.MulticastResponse) *MulticastResponse

type MulticastResponse struct {
	Successful, Failed MulticastStateResponse
}

type MulticastStateResponse struct {
	ServerIds []int64
	Responses []interface{}
}

func NewMulticastController(multicaster gateway.Multicaster, responseTransformer multicastResponseTransformer) *MulticastController {
	return &MulticastController{multicaster, responseTransformer}
}

type MulticastRequest struct {
	Call          string
	HubSessionKey string
	ServerIDs     []int64
	ArgsByServer  map[int64][]interface{}
}

func (h *MulticastController) Multicast(r *http.Request, args *MulticastRequest, reply *struct{ Data *MulticastResponse }) error {
	multicastResponse, err := h.multicaster.Multicast(args.HubSessionKey, args.Call, args.ServerIDs, args.ArgsByServer)
	if err != nil {
		return err
	}
	reply.Data = h.responseTransformer(multicastResponse)
	return nil
}
