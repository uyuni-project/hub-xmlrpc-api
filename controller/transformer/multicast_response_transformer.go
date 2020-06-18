package transformer

import (
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

// MulticastResponseTransformer turns a multicast response from the gateway to the controller format
func MulticastResponseTransformer(multicastResponse *gateway.MulticastResponse) *controller.MulticastResponse {
	return &controller.MulticastResponse{
		transformToSuccessfulResponses(multicastResponse.SuccessfulResponses),
		transformToFailedResponses(multicastResponse.FailedResponses),
	}
}

func transformToSuccessfulResponses(serverCallResponses map[int64]gateway.ServerSuccessfulResponse) controller.MulticastStateResponse {
	serverIDs := make([]int64, 0, len(serverCallResponses))
	responses := make([]interface{}, 0, len(serverCallResponses))

	for serverID, response := range serverCallResponses {
		serverIDs = append(serverIDs, serverID)
		responses = append(responses, response.Response)
	}
	return controller.MulticastStateResponse{serverIDs, responses}
}

func transformToFailedResponses(serverCallResponses map[int64]gateway.ServerFailedResponse) controller.MulticastStateResponse {
	serverIDs := make([]int64, 0, len(serverCallResponses))
	responses := make([]interface{}, 0, len(serverCallResponses))

	for serverID, response := range serverCallResponses {
		serverIDs = append(serverIDs, serverID)
		responses = append(responses, response.ErrorMessage)
	}
	return controller.MulticastStateResponse{serverIDs, responses}
}
