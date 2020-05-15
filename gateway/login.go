package gateway

import (
	"log"
)

const (
	manualLoginMode = iota // 0
	relayLoginMode         // 1
)

//HubLoginer interface for Login operations
type HubLoginer interface {
	Login(username, password string) (string, error)
	LoginWithAuthRelayMode(username, password string) (string, error)
	LoginWithAutoconnectMode(username, password string) (string, error)
}

type hubLoginer struct {
	hubAPIEndpoint                   string
	uyuniAuthenticator               UyuniAuthenticator
	serverAuthenticator              ServerAuthenticator
	uyuniServerTopologyInfoRetriever UyuniTopologyInfoRetriever
	hubSessionRepository             HubSessionRepository
}

//NewHubLoginer instantiates a hubLoginer
func NewHubLoginer(hubAPIEndpoint string, uyuniAuthenticator UyuniAuthenticator,
	serverAuthenticator ServerAuthenticator, uyuniTopologyInfoRetriever UyuniTopologyInfoRetriever,
	hubSessionRepository HubSessionRepository) *hubLoginer {
	return &hubLoginer{hubAPIEndpoint, uyuniAuthenticator, serverAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository}
}

func (h *hubLoginer) Login(username, password string) (string, error) {
	return h.loginToHub(username, password, manualLoginMode)
}

func (h *hubLoginer) LoginWithAuthRelayMode(username, password string) (string, error) {
	return h.loginToHub(username, password, relayLoginMode)
}

func (h *hubLoginer) LoginWithAutoconnectMode(username, password string) (string, error) {
	hubSessionKey, err := h.LoginWithAuthRelayMode(username, password)
	if err != nil {
		return "", err
	}
	userServerIDs, err := h.uyuniServerTopologyInfoRetriever.RetrieveUserServerIDs(h.hubAPIEndpoint, hubSessionKey, username)
	if err != nil {
		return "", err
	}
	h.serverAuthenticator.AttachToServers(hubSessionKey, userServerIDs, nil)
	return hubSessionKey, nil
}

func (h *hubLoginer) loginToHub(username, password string, loginMode int) (string, error) {
	hubToken, err := h.uyuniAuthenticator.Login(h.hubAPIEndpoint, username, password)
	if err != nil {
		log.Printf("Error ocurred while trying to login into the Hub: %v", err)
		return "", err
	}
	h.hubSessionRepository.SaveHubSession(NewHubSession(hubToken, username, password, loginMode))
	return hubToken, nil
}
