package gateway

import (
	"errors"
	"log"
)

const (
	loginPath           = "auth.login"
	listUserSystemsPath = "system.listUserSystems"
	listSystemFQDNsPath = "system.listFqdns"
	systemIDField       = "id"

	manualLoginMode      = iota // 0
	relayLoginMode              // 1
	autoconnectLoginMode        // 2
)

type Authenticator interface {
	Login(username, password string) (string, error)
	LoginWithAuthRelayMode(username, password string) (string, error)
	LoginWithAutoconnectMode(username, password string) (string, error)
	AttachToServers(hubSessionKey string, argsByServerID map[int64][]interface{}) (*MulticastResponse, error)
}

type authenticator struct {
	client         Client
	session        Session
	hubAPIEndpoint string
}

func NewAuthenticator(client Client, session Session, hubAPIEndpoint string) *authenticator {
	return &authenticator{client, session, hubAPIEndpoint}
}

func (a *authenticator) Login(username, password string) (string, error) {
	hubSessionKey, err := a.loginToHub(username, password, manualLoginMode)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return hubSessionKey, nil
}

func (h *authenticator) LoginWithAuthRelayMode(username, password string) (string, error) {
	hubSessionKey, err := h.loginToHub(username, password, relayLoginMode)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return hubSessionKey, nil
}

func (a *authenticator) LoginWithAutoconnectMode(username, password string) (string, error) {
	hubSessionKey, err := a.loginToHub(username, password, autoconnectLoginMode)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}

	userServerIDs, err := a.getUserServerIDs(hubSessionKey, username)
	if err != nil {
		//TODO: should we return an error? retry the login or what?
	}
	_, err = a.loginIntoServersUsingSameCredentials(userServerIDs, username, password, hubSessionKey)
	if err != nil {
		//TODO: should we return an error? retry the login or what?
	}
	return hubSessionKey, nil
}

func (a *authenticator) AttachToServers(hubSessionKey string, credentialsByServer map[int64][]interface{}) (*MulticastResponse, error) {
	hubSession := a.session.RetrieveHubSession(hubSessionKey)
	if hubSession == nil {
		log.Printf("HubSession was not found: %v", hubSessionKey)
		return nil, errors.New("Authentication error: provided session key is invalid")
	}

	if hubSession.loginMode == relayLoginMode {
		serverIDs := make([]int64, 0)
		for serverID := range credentialsByServer {
			serverIDs = append(serverIDs, serverID)
		}
		return a.loginIntoServersUsingSameCredentials(serverIDs, hubSession.username, hubSession.password, hubSessionKey)
	}
	return a.loginIntoSystems(hubSessionKey, credentialsByServer)
}

func (a *authenticator) loginIntoServersUsingSameCredentials(serverIDs []int64, username, password, hubSessionKey string) (*MulticastResponse, error) {
	credentialsByServer := generateSameCredentialsForServers(serverIDs, username, password)
	return a.loginIntoSystems(hubSessionKey, credentialsByServer)
}

func (a *authenticator) loginToHub(username, password string, loginMode int) (string, error) {
	response, err := a.client.ExecuteCall(a.hubAPIEndpoint, loginPath, []interface{}{username, password})
	if err != nil {
		log.Printf("Error ocurred while trying to login into the Hub: %v", err)
		return "", err
	}
	hubSessionKey := response.(string)
	a.session.SaveHubSession(NewHubSession(hubSessionKey, username, password, loginMode))
	return hubSessionKey, nil
}

func (a *authenticator) getUserServerIDs(hubSessionKey, username string) ([]int64, error) {
	userServers, err := a.client.ExecuteCall(a.hubAPIEndpoint, listUserSystemsPath, []interface{}{hubSessionKey, username})
	if err != nil {
		log.Printf("Error ocurred while trying to login into the user systems: %v", err)
		return nil, err
	}
	userServersSlice := userServers.([]interface{})

	serverIDs := make([]int64, 0, len(userServersSlice))
	for _, userSystem := range userServersSlice {
		serverID := userSystem.(map[string]interface{})[systemIDField].(int64)
		serverIDs = append(serverIDs, serverID)
	}
	return serverIDs, nil
}

func (a *authenticator) loginIntoSystems(hubSessionKey string, credentialsByServerID map[int64][]interface{}) (*MulticastResponse, error) {
	loginIntoSystemsArgs, err := a.generateServerCallRequests(hubSessionKey, credentialsByServerID)
	if err != nil {
		//TODO: what to do with the error here?
	}
	multicastResponse := executeCallOnServers(loginIntoSystemsArgs, a.client)
	a.saveServerSessions(hubSessionKey, multicastResponse)
	return multicastResponse, nil
}

func (a *authenticator) saveServerSessions(hubSessionKey string, loginResponses *MulticastResponse) {
	serverSessions := make(map[int64]*ServerSession)
	for _, response := range loginResponses.SuccessfulResponses {
		serverSessions[response.ServerID] = &ServerSession{response.ServerID, response.endpoint, response.Response.(string), hubSessionKey}
	}
	// TODO: If we don't save responses for failed servers in session, user will get `Invalid session error" because of failed lookup later
	// and wouldn't even get results for those where call was successful. We need a better mechanism to handle such cases.
	//save for failed as well
	for _, response := range loginResponses.FailedResponses {
		serverSessions[response.ServerID] = &ServerSession{response.ServerID, response.endpoint, "login-error", hubSessionKey}
	}
	a.session.SaveServerSessions(hubSessionKey, serverSessions)
}

func (a *authenticator) generateServerCallRequests(hubSessionKey string, credentialsByServer map[int64][]interface{}) ([]serverCallRequest, error) {
	serverCallRequests := make([]serverCallRequest, 0, len(credentialsByServer))

	for serverID, credentials := range credentialsByServer {
		serverAPIEndpoint, err := a.retrieveServerAPIEndpoint(hubSessionKey, serverID)
		if err != nil {
			//TODO: what to do with failing servers?
		} else {
			serverCallRequests = append(serverCallRequests, serverCallRequest{serverAPIEndpoint, loginPath, serverID, credentials})
		}
	}
	return serverCallRequests, nil
}

func (a *authenticator) retrieveServerAPIEndpoint(hubSessionKey string, serverID int64) (string, error) {
	//TODO: we should deal with cases when we have more than one fqdn
	response, err := a.client.ExecuteCall(a.hubAPIEndpoint, listSystemFQDNsPath, []interface{}{hubSessionKey, serverID})
	if err != nil {
		log.Printf("Error ocurred when retrieving the system Fqdns for serverID: %v, error:%v", serverID, err)
		return "", err
	}
	//TODO: check for casting errors.
	//TODO: check the fqdn array is not empty
	firstFqdn := response.([]interface{})[0].(string)
	return "http://" + firstFqdn + "/rpc/api", nil
}

func generateSameCredentialsForServers(serverIDs []int64, username, password string) map[int64][]interface{} {
	credentialsByServerID := make(map[int64][]interface{})
	for _, serverID := range serverIDs {
		credentialsByServerID[serverID] = []interface{}{username, password}
	}
	return credentialsByServerID
}
