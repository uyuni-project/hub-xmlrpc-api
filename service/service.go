package service

import "log"

type service struct {
	client        Client
	session       Session
	hubSumaAPIURL string
}

func (s *service) isHubSessionValid(hubSessionKey string) bool {
	isValid, err := s.client.ExecuteCall(s.hubSumaAPIURL, "auth.isSessionKeyValid", []interface{}{hubSessionKey})
	if err != nil {
		log.Printf("Login error: %v", err)
		s.session.RemoveHubSession(hubSessionKey)
		return false
	}
	return isValid.(bool)
}
