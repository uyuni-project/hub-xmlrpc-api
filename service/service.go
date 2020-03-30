package service

import "log"

const CHECK_HUB_SESSION_KEY_PATH = "auth.isSessionKeyValid"

type service struct {
	client        Client
	session       Session
	hubSumaAPIURL string
}

func (s *service) isHubSessionValid(hubSessionKey string) bool {
	isValid, err := s.client.ExecuteCall(s.hubSumaAPIURL, CHECK_HUB_SESSION_KEY_PATH, []interface{}{hubSessionKey})
	if err != nil {
		log.Printf("Login error: %v", err)
		s.session.RemoveHubSession(hubSessionKey)
		return false
	}
	return isValid.(bool)
}
