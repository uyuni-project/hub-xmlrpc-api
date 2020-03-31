package gateway

import (
	"errors"
	"log"
)

const (
	loginPath             = "auth.login"
	listUserSystemsPath   = "system.listUserSystems"
	listSystemFQDNsPath   = "system.listFqdns"
	isSessionKeyValidPath = "auth.isSessionKeyValid"

	manualLoginMode      = iota // 0
	relayLoginMode              // 1
	autoconnectLoginMode        // 2
)

type Authorizer interface {
	Login(username, password string) (string, error)
	LoginWithAuthRelayMode(username, password string) (string, error)
	LoginWithAutoconnectMode(username, password string) (string, error)
	AttachToServers(hubSessionKey string, argsByServerID map[int64][]interface{}) (*MulticastResponse, error)
}

type sessionValidator interface {
	isHubSessionValid(hubSessionKey string) bool
}

type AuthorizationService struct {
	client        Client
	session       Session
	hubSumaAPIURL string
}

func NewAuthorizationService(client Client, session Session, hubSumaAPIURL string) *AuthorizationService {
	return &AuthorizationService{client: client, session: session, hubSumaAPIURL: hubSumaAPIURL}
}

func (a *AuthorizationService) Login(username, password string) (string, error) {
	hubSessionKey, err := a.loginToHub(username, password, manualLoginMode)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return hubSessionKey, nil
}

func (h *AuthorizationService) LoginWithAuthRelayMode(username, password string) (string, error) {
	hubSessionKey, err := h.loginToHub(username, password, relayLoginMode)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return hubSessionKey, nil
}

func (a *AuthorizationService) LoginWithAutoconnectMode(username, password string) (string, error) {
	hubSessionKey, err := a.loginToHub(username, password, autoconnectLoginMode)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}

	err = a.loginIntoUserSystems(hubSessionKey, username, password)
	if err != nil {
		//TODO: should we return an error? retry the login or what?
	}
	return hubSessionKey, nil
}

func (a *AuthorizationService) AttachToServers(hubSessionKey string, argsByServerID map[int64][]interface{}) (*MulticastResponse, error) {
	if a.isHubSessionValid(hubSessionKey) {

		hubSession := a.session.RetrieveHubSession(hubSessionKey)
		if hubSession == nil {
			log.Printf("HubSessionKey was not found: %v", hubSessionKey)
			//TODO: what error should we return here?
			return nil, errors.New("provided session key is invalid")
		}

		credentialsByServerID := argsByServerID
		if hubSession.loginMode == relayLoginMode {

			credentialsByServerID = make(map[int64][]interface{})
			for serverID := range argsByServerID {
				credentialsByServerID[serverID] = []interface{}{hubSession.username, hubSession.password}
			}
		}
		return a.loginIntoSystems(hubSessionKey, credentialsByServerID)
	}
	log.Printf("Provided session key is invalid: %v", hubSessionKey)
	//TODO: should we return an error here?
	return nil, nil
}

func (a *AuthorizationService) loginToHub(username, password string, loginMode int) (string, error) {
	response, err := a.client.ExecuteCall(a.hubSumaAPIURL, loginPath, []interface{}{username, password})
	if err != nil {
		log.Printf("Error ocurred while trying to login into the Hub: %v", err)
		return "", err
	}
	hubSessionKey := response.(string)
	a.session.SaveHubSession(hubSessionKey, &HubSession{username, password, loginMode})
	return hubSessionKey, nil
}

func (a *AuthorizationService) loginIntoUserSystems(hubSessionKey, username, password string) error {
	userSystems, err := a.client.ExecuteCall(a.hubSumaAPIURL, listUserSystemsPath, []interface{}{hubSessionKey, username})
	if err != nil {
		log.Printf("Error ocurred while trying to login into the user systems: %v", err)
		return err
	}
	userSystemsSlice := userSystems.([]interface{})

	credentialsByServerID := make(map[int64][]interface{})
	for _, userSystem := range userSystemsSlice {
		serverID := userSystem.(map[string]interface{})["id"].(int64)
		credentialsByServerID[serverID] = []interface{}{username, password}
	}

	//TODO: what to do with the response here?
	_, err = a.loginIntoSystems(hubSessionKey, credentialsByServerID)
	return err
}

func (a *AuthorizationService) loginIntoSystems(hubSessionKey string, credentialsByServerID map[int64][]interface{}) (*MulticastResponse, error) {
	loginIntoSystemsArgs, serverURLByServerID, err := a.resolveLoginIntoSystemsArgs(hubSessionKey, credentialsByServerID)
	if err != nil {
		//TODO: what to do with the error here?
	}
	responses := executeMulticastCall(loginPath, loginIntoSystemsArgs, a.client)
	successfulResponses := responses.Successful
	failedResponses := responses.Failed

	//save in session
	for i, serverID := range successfulResponses.ServerIds {
		a.session.SaveServerSession(hubSessionKey, serverID, &ServerSession{serverURLByServerID[serverID], successfulResponses.Responses[i].(string)})
	}

	// TODO: If we don't save responses for failed servers in session, user will get `Invalid session error" because of failed lookup later
	// and wouldn't even get results for those where call was successful. We need a better mechanism to handle such cases.

	//save for failed as well
	for _, serverID := range failedResponses.ServerIds {
		a.session.SaveServerSession(hubSessionKey, serverID, &ServerSession{serverURLByServerID[serverID], "login-error"})
	}
	return responses, nil
}

func (a *AuthorizationService) resolveLoginIntoSystemsArgs(hubSessionKey string, credentialsByServerID map[int64][]interface{}) ([]multicastServerArgs, map[int64]string, error) {
	multicastArgs := make([]multicastServerArgs, 0, len(credentialsByServerID))
	serverURLByServerID := make(map[int64]string)

	for serverID, credentials := range credentialsByServerID {
		url, err := a.retrieveServerAPIURL(hubSessionKey, serverID)
		if err != nil {
			//TODO: what to do with failing servers?
		} else {
			serverURLByServerID[serverID] = url
			multicastArgs = append(multicastArgs, multicastServerArgs{url, serverID, credentials})
		}
	}
	return multicastArgs, serverURLByServerID, nil
}

func (a *AuthorizationService) retrieveServerAPIURL(hubSessionKey string, serverID int64) (string, error) {
	//TODO: we should deal with cases when we have more than one fqdn
	response, err := a.client.ExecuteCall(a.hubSumaAPIURL, listSystemFQDNsPath, []interface{}{hubSessionKey, serverID})
	if err != nil {
		log.Printf("Error ocurred when retrieving the system Fqdns for serverID: %v, error:%v", serverID, err)
		return "", err
	}
	//TODO: check for casting errors.
	//TODO: check the fqdn array is not empty
	firstFqdn := response.([]interface{})[0].(string)
	return "http://" + firstFqdn + "/rpc/api", nil
}

func (a *AuthorizationService) isHubSessionValid(hubSessionKey string) bool {
	isValid, err := a.client.ExecuteCall(a.hubSumaAPIURL, isSessionKeyValidPath, []interface{}{hubSessionKey})
	if err != nil {
		log.Printf("Login error: %v", err)
		a.session.RemoveHubSession(hubSessionKey)
		return false
	}
	return isValid.(bool)
}
