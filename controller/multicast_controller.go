package controller

import (
	"net/http"
	"strings"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

type MulticastController struct {
	multicaster gateway.Multicaster
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

func (h *MulticastController) Multicast(r *http.Request, args *gateway.MulticastRequest, reply *struct{ Data *MulticastResponse }) error {
	multicastResponse, err := h.multicaster.Multicast(args)
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
	return &MulticastResponse{
		transformToSuccessfulResponses(multicastResponse.SuccessfulResponses),
		transformToFailedResponses(multicastResponse.FailedResponses),
	}
}

func transformToSuccessfulResponses(serverCallResponses []*gateway.ServerSuccessfulResponse) MulticastStateResponse {
	serverIDs := make([]int64, 0, len(serverCallResponses))
	responses := make([]interface{}, 0, len(serverCallResponses))

	for _, response := range serverCallResponses {
		serverIDs = append(serverIDs, response.ServerID)
		responses = append(responses, response.Response)
	}
	return MulticastStateResponse{serverIDs, responses}
}

func transformToFailedResponses(serverCallResponses []*gateway.ServerFailedResponse) MulticastStateResponse {
	serverIDs := make([]int64, 0, len(serverCallResponses))
	responses := make([]interface{}, 0, len(serverCallResponses))

	for _, response := range serverCallResponses {
		serverIDs = append(serverIDs, response.ServerID)
		responses = append(responses, response.ErrorMessage)
	}
	return MulticastStateResponse{serverIDs, responses}
}
