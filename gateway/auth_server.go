package gateway

import (
	"errors"
	"log"
)

type ServerAuthenticator interface {
	AttachToServers(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error)
	loginToServersUsingSameCredentials(serverIDs []int64, username, password, hubSessionKey string) (*MulticastResponse, error)
}

type Credentials struct {
	Username, Password string
}

type serverAuthenticator struct {
	uyuniServerAuthenticator      UyuniServerAuthenticator
	uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever
	hubSessionRepository          HubSessionRepository
	serverSessionRepository       ServerSessionRepository
}

func NewServerAuthenticator(uyuniServerAuthenticator UyuniServerAuthenticator,
	uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever, hubSessionRepository HubSessionRepository, serverSessionRepository ServerSessionRepository) *serverAuthenticator {
	return &serverAuthenticator{uyuniServerAuthenticator, uyuniHubTopologyInfoRetriever, hubSessionRepository, serverSessionRepository}
}

func (a *serverAuthenticator) AttachToServers(hubSessionKey string, serverIDs []int64, credentialsByServer map[int64]*Credentials) (*MulticastResponse, error) {
	hubSession := a.hubSessionRepository.RetrieveHubSession(hubSessionKey)
	if hubSession == nil {
		log.Printf("HubSession was not found: %v", hubSessionKey)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	if hubSession.loginMode == relayLoginMode {
		return a.loginToServersUsingSameCredentials(serverIDs, hubSession.username, hubSession.password, hubSessionKey)
	}
	return a.loginToServers(serverIDs, credentialsByServer, hubSessionKey)
}

func (a *serverAuthenticator) loginToServersUsingSameCredentials(serverIDs []int64, username, password, hubSessionKey string) (*MulticastResponse, error) {
	credentialsByServer := generateSameCredentialsForServers(serverIDs, username, password)
	return a.loginToServers(serverIDs, credentialsByServer, hubSessionKey)
}

func (a *serverAuthenticator) loginToServers(serverIDs []int64, credentialsByServer map[int64]*Credentials, hubSessionKey string) (*MulticastResponse, error) {
	endpointByServer, err := a.uyuniHubTopologyInfoRetriever.RetrieveServerAPIEndpoints(hubSessionKey, serverIDs)
	if err != nil {
		//TODO: what to do with failing servers?
	}
	multicastCallRequest, err := a.generateLoginMuticastCallRequest(credentialsByServer, endpointByServer)
	if err != nil {
		//TODO: what to do with the error here?
	}
	loginResponse := executeCallOnServers(multicastCallRequest)
	a.saveServerSessions(hubSessionKey, loginResponse)
	return loginResponse, nil
}

func (a *serverAuthenticator) generateLoginMuticastCallRequest(credentialsByServer map[int64]*Credentials, endpointByServer map[int64]string) (*multicastCallRequest, error) {
	call := func(endpoint string, args []interface{}) (interface{}, error) {
		return a.uyuniServerAuthenticator.Login(endpoint, args[0].(string), args[1].(string))
	}

	serverCallInfos := make([]serverCallInfo, 0, len(credentialsByServer))
	for serverID, endpoint := range endpointByServer {
		args := []interface{}{credentialsByServer[serverID].Username, credentialsByServer[serverID].Password}
		serverCallInfos = append(serverCallInfos, serverCallInfo{serverID, endpoint, args})
	}
	return &multicastCallRequest{call, serverCallInfos}, nil
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
