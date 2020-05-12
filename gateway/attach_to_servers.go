package gateway

import (
	"errors"
	"log"
)

type ServerAuthenticator interface {
	AttachToServers(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error)
	attachServersToHubSessionUsingSameCredentials(serverIDs []int64, username, password, hubSessionKey string) (*MulticastResponse, error)
}

type Credentials struct {
	Username, Password string
}

type serverAuthenticator struct {
	hubAPIEndpoint             string
	uyuniAuthenticator         UyuniAuthenticator
	uyuniTopologyInfoRetriever UyuniTopologyInfoRetriever
	hubSessionRepository       HubSessionRepository
	serverSessionRepository    ServerSessionRepository
}

func NewServerAuthenticator(hubAPIEndpoint string, uyuniAuthenticator UyuniAuthenticator,
	uyuniTopologyInfoRetriever UyuniTopologyInfoRetriever, hubSessionRepository HubSessionRepository,
	serverSessionRepository ServerSessionRepository) *serverAuthenticator {
	return &serverAuthenticator{hubAPIEndpoint, uyuniAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository, serverSessionRepository}
}

func (a *serverAuthenticator) AttachToServers(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error) {
	hubSession := a.hubSessionRepository.RetrieveHubSession(hubSessionKey)
	if hubSession == nil {
		log.Printf("HubSession was not found: %v", hubSessionKey)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	if hubSession.loginMode == relayLoginMode {
		return a.attachServersToHubSessionUsingSameCredentials(serverIDs, hubSession.username, hubSession.password, hubSessionKey)
	}
	return a.attachServersToHubSession(serverIDs, credentialsByServer, hubSessionKey)
}

func (a *serverAuthenticator) attachServersToHubSessionUsingSameCredentials(serverIDs []int64, username, password, hubSessionKey string) (*MulticastResponse, error) {
	credentialsByServer := generateSameCredentialsForServers(serverIDs, username, password)
	return a.attachServersToHubSession(serverIDs, credentialsByServer, hubSessionKey)
}

func (a *serverAuthenticator) attachServersToHubSession(serverIDs []int64, credentialsByServer map[int64]*Credentials, hubSessionKey string) (*MulticastResponse, error) {
	endpointByServer, err := a.uyuniTopologyInfoRetriever.RetrieveServerAPIEndpoints(a.hubAPIEndpoint, hubSessionKey, serverIDs)
	if err != nil {
		//TODO: what to do with failing servers?
	}
	multicastCallRequest := a.generateLoginMuticastCallRequest(credentialsByServer, endpointByServer)
	loginResponse := executeCallOnServers(multicastCallRequest)
	a.saveServerSessions(hubSessionKey, loginResponse)
	return loginResponse, nil
}

func (a *serverAuthenticator) generateLoginMuticastCallRequest(credentialsByServer map[int64]*Credentials, endpointByServer map[int64]string) *multicastCallRequest {
	call := func(endpoint string, args []interface{}) (interface{}, error) {
		return a.uyuniAuthenticator.Login(endpoint, args[0].(string), args[1].(string))
	}
	serverCallInfos := make([]serverCallInfo, 0, len(credentialsByServer))
	for serverID, endpoint := range endpointByServer {
		args := []interface{}{credentialsByServer[serverID].Username, credentialsByServer[serverID].Password}
		serverCallInfos = append(serverCallInfos, serverCallInfo{serverID, endpoint, args})
	}
	return &multicastCallRequest{call, serverCallInfos}
}

func (a *serverAuthenticator) saveServerSessions(hubSessionKey string, loginResponses *MulticastResponse) {
	serverSessions := make(map[int64]*ServerSession)
	for serverID, response := range loginResponses.SuccessfulResponses {
		serverSessions[serverID] = &ServerSession{serverID, response.endpoint, response.Response.(string), hubSessionKey}
	}
	// TODO:
	//save for failed as well
	for serverID, response := range loginResponses.FailedResponses {
		serverSessions[serverID] = &ServerSession{serverID, response.endpoint, "login-error", hubSessionKey}
	}
	a.serverSessionRepository.SaveServerSessions(hubSessionKey, serverSessions)
}

func generateSameCredentialsForServers(serverIDs []int64, username, password string) map[int64]*Credentials {
	credentialsByServerID := make(map[int64]*Credentials)
	for _, serverID := range serverIDs {
		credentialsByServerID[serverID] = &Credentials{username, password}
	}
	return credentialsByServerID
}
