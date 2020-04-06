package gateway

import (
	"log"
)

const (
	manualLoginMode      = iota // 0
	relayLoginMode              // 1
	autoconnectLoginMode        // 2
)

type HubAuthenticator interface {
	Login(username, password string) (string, error)
	LoginWithAuthRelayMode(username, password string) (string, error)
	LoginWithAutoconnectMode(username, password string) (string, error)
}

type hubAuthenticator struct {
	uyuniHubAuthenticator         UyuniHubAuthenticator
	serverAuthenticator           ServerAuthenticator
	uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever
	session                       Session
}

func NewHubAuthenticator(uyuniHubAuthenticator UyuniHubAuthenticator, serverAuthenticator ServerAuthenticator, uyuniHubTopologyInfoRetriever UyuniHubTopologyInfoRetriever, session Session) *hubAuthenticator {
	return &hubAuthenticator{uyuniHubAuthenticator, serverAuthenticator, uyuniHubTopologyInfoRetriever, session}
}

func (h *hubAuthenticator) Login(username, password string) (string, error) {
	return h.loginToHub(username, password, manualLoginMode)
}

func (h *hubAuthenticator) LoginWithAuthRelayMode(username, password string) (string, error) {
	return h.loginToHub(username, password, relayLoginMode)
}

func (h *hubAuthenticator) LoginWithAutoconnectMode(username, password string) (string, error) {
	hubSessionKey, err := h.loginToHub(username, password, autoconnectLoginMode)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	userServerIDs, err := h.uyuniHubTopologyInfoRetriever.RetrieveUserServerIDs(hubSessionKey, username)
	if err != nil {
		//TODO: should we return an error? retry the login or what?
	}
	h.serverAuthenticator.loginToServersUsingSameCredentials(userServerIDs, username, password, hubSessionKey)
	return hubSessionKey, nil
}

func (h *hubAuthenticator) loginToHub(username, password string, loginMode int) (string, error) {
	hubToken, err := h.uyuniHubAuthenticator.Login(username, password)
	if err != nil {
		log.Printf("Error ocurred while trying to login into the Hub: %v", err)
		return "", err
	}
	h.session.SaveHubSession(NewHubSession(hubToken, username, password, loginMode))
	return hubToken, nil
}
