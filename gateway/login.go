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
	hubAPIEndpoint                   string
	uyuniAuthenticator               UyuniAuthenticator
	serverAuthenticator              ServerAuthenticator
	uyuniServerTopologyInfoRetriever UyuniTopologyInfoRetriever
	hubSessionRepository             HubSessionRepository
}

func NewHubAuthenticator(hubAPIEndpoint string, uyuniAuthenticator UyuniAuthenticator,
	serverAuthenticator ServerAuthenticator, uyuniTopologyInfoRetriever UyuniTopologyInfoRetriever,
	hubSessionRepository HubSessionRepository) *hubAuthenticator {
	return &hubAuthenticator{hubAPIEndpoint, uyuniAuthenticator, serverAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository}
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
		return "", err
	}
	userServerIDs, err := h.uyuniServerTopologyInfoRetriever.RetrieveUserServerIDs(h.hubAPIEndpoint, hubSessionKey, username)
	if err != nil {
		return "", err
	}
	h.serverAuthenticator.attachServersToHubSessionUsingSameCredentials(userServerIDs, username, password, hubSessionKey)
	return hubSessionKey, nil
}

func (h *hubAuthenticator) loginToHub(username, password string, loginMode int) (string, error) {
	hubToken, err := h.uyuniAuthenticator.Login(h.hubAPIEndpoint, username, password)
	if err != nil {
		log.Printf("Error ocurred while trying to login into the Hub: %v", err)
		return "", err
	}
	h.hubSessionRepository.SaveHubSession(NewHubSession(hubToken, username, password, loginMode))
	return hubToken, nil
}
