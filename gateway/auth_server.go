package gateway

import (
	"errors"
	"log"
)

type ServerAuthenticator interface {
	AttachToServers(request *AttachToServersRequest) (*MulticastResponse, error)
	loginToServersUsingSameCredentials(serverIDs []int64, username, password, hubSessionKey string) (*MulticastResponse, error)
}

type AttachToServersRequest struct {
	HubSessionKey       string
	ServerIDs           []int64
	CredentialsByServer map[int64]*Credentials
}

type Credentials struct {
	Username, Password string
}

type serverAuthenticator struct {
	uyuniServerAuthenticator      UyuniServerAuthenticator
	uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever
	session                       Session
}

func NewServerAuthenticator(uyuniServerAuthenticator UyuniServerAuthenticator,
	uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever, session Session) *serverAuthenticator {
	return &serverAuthenticator{uyuniServerAuthenticator, uyuniHubTopologyInfoRetriever, session}
}

func (a *serverAuthenticator) AttachToServers(request *AttachToServersRequest) (*MulticastResponse, error) {
	hubSession := a.session.RetrieveHubSession(request.HubSessionKey)
	if hubSession == nil {
		log.Printf("HubSession was not found: %v", request.HubSessionKey)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}
	if hubSession.loginMode == relayLoginMode {
		return a.loginToServersUsingSameCredentials(request.ServerIDs, hubSession.username, hubSession.password, request.HubSessionKey)
	}
	return a.loginToServers(request.ServerIDs, request.CredentialsByServer, request.HubSessionKey)
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
	// TODO: If we don't save responses for failed servers in session, user will get `Invalid session error" because of failed lookup later
	// and wouldn't even get results for those where call was successful. We need a better mechanism to handle such cases.
	//save for failed as well
	for serverID, response := range loginResponses.FailedResponses {
		serverSessions[serverID] = &ServerSession{serverID, response.endpoint, "login-error", hubSessionKey}
	}
	a.session.SaveServerSessions(hubSessionKey, serverSessions)
}

func generateSameCredentialsForServers(serverIDs []int64, username, password string) map[int64]*Credentials {
	credentialsByServerID := make(map[int64]*Credentials)
	for _, serverID := range serverIDs {
		credentialsByServerID[serverID] = &Credentials{username, password}
	}
	return credentialsByServerID
}
