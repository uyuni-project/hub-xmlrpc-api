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

type MulticastResponse struct {
	Successful, Failed MulticastStateResponse
}

type MulticastStateResponse struct {
	ServerIds []int64
	Responses []interface{}
}

func NewMulticastController(multicaster gateway.Multicaster) *MulticastController {
	return &MulticastController{multicaster}
}

func (h *MulticastController) Multicast(r *http.Request, args *MulticastRequest, reply *struct{ Data *MulticastResponse }) error {
	method := removeMulticastNamespace(args.Method)
	multicastResponse, err := h.multicaster.Multicast(args.HubSessionKey, method, args.ArgsByServer)
	if err != nil {
		return err
	}
	reply.Data = transformToOutputModel(multicastResponse)
	return nil
}

func removeMulticastNamespace(method string) string {
	parts := strings.Split(method, ".")
	slice := parts[1:len(parts)]
	return strings.Join(slice, ".")
}

func transformToOutputModel(multicastResponse *gateway.MulticastResponse) *MulticastResponse {
	successfulKeys, successfulValues := getServerIDsAndResponses(multicastResponse.SuccessfulResponses)
	failedKeys, failedValues := getServerIDsAndResponses(multicastResponse.FailedResponses)

	return &MulticastResponse{MulticastStateResponse{successfulKeys, successfulValues}, MulticastStateResponse{failedKeys, failedValues}}
}

func getServerIDsAndResponses(in map[int64]interface{}) ([]int64, []interface{}) {
	serverIDs := make([]int64, 0, len(in))
	responses := make([]interface{}, 0, len(in))

	for serverID, response := range in {
		serverIDs = append(serverIDs, serverID)
		responses = append(responses, response)
	}
	return serverIDs, responses
}
