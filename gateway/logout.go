package gateway

import (
	"errors"
	"log"
)

func (h *hubAuthenticator) Logout(hubSessionKey string) error {
	hubSession := h.hubSessionRepository.RetrieveHubSession(hubSessionKey)
	if hubSession == nil {
		log.Printf("HubSession was not found. HubSessionKey: %v", hubSessionKey)
		return errors.New("Authentication error: provided session key is invalid")
	}
	err := h.uyuniAuthenticator.Logout(h.hubAPIEndpoint, hubSessionKey)
	if err != nil {
		return err
	}
	h.logoutFromServersInHubSession(hubSession.ServerSessions)
	h.hubSessionRepository.RemoveHubSession(hubSessionKey)
	return nil
}

func (h *hubAuthenticator) logoutFromServersInHubSession(serverSessions map[int64]*ServerSession) *MulticastResponse {
	multicastCallRequest := h.generateLogoutMuticastCallRequest(serverSessions)
	return executeCallOnServers(multicastCallRequest)
}

func (h *hubAuthenticator) generateLogoutMuticastCallRequest(serverSessions map[int64]*ServerSession) *multicastCallRequest {
	call := func(endpoint string, args []interface{}) (interface{}, error) {
		return nil, h.uyuniAuthenticator.Logout(endpoint, args[0].(string))
	}
	serverCallInfos := make([]serverCallInfo, 0, len(serverSessions))
	for serverID, serverSession := range serverSessions {
		serverCallInfos = append(serverCallInfos, serverCallInfo{serverID, serverSession.serverAPIEndpoint, []interface{}{serverSession.serverSessionKey}})
	}
	return &multicastCallRequest{call, serverCallInfos}
}
