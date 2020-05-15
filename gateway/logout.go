package gateway

import (
	"errors"
	"log"
)

//HubLogouter provides an interface for logout operations
type HubLogouter interface {
	Logout(hubSessionKey string) error
}

type hubLogouter struct {
	hubAPIEndpoint       string
	uyuniAuthenticator   UyuniAuthenticator
	hubSessionRepository HubSessionRepository
}

//NewHubLogouter instantiates a HubLogouter
func NewHubLogouter(hubAPIEndpoint string, uyuniAuthenticator UyuniAuthenticator, hubSessionRepository HubSessionRepository) *hubLogouter {
	return &hubLogouter{hubAPIEndpoint, uyuniAuthenticator, hubSessionRepository}
}

func (h *hubLogouter) Logout(hubSessionKey string) error {
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

func (h *hubLogouter) logoutFromServersInHubSession(serverSessions map[int64]*ServerSession) *MulticastResponse {
	multicastCallRequest := h.generateLogoutMuticastCallRequest(serverSessions)
	return executeCallOnServers(multicastCallRequest)
}

func (h *hubLogouter) generateLogoutMuticastCallRequest(serverSessions map[int64]*ServerSession) *multicastCallRequest {
	call := func(endpoint string, args []interface{}) (interface{}, error) {
		return nil, h.uyuniAuthenticator.Logout(endpoint, args[0].(string))
	}
	serverCallInfos := make([]serverCallInfo, 0, len(serverSessions))
	for serverID, serverSession := range serverSessions {
		serverCallInfos = append(serverCallInfos, serverCallInfo{serverID, serverSession.serverAPIEndpoint, []interface{}{serverSession.serverSessionKey}})
	}
	return &multicastCallRequest{call, serverCallInfos}
}
