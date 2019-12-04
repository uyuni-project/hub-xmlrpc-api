package session

type ApiSession struct {
	hubSessionKey, username, password string
	userServerURLByKey                map[string]string
}

//TODO: remove when Abids PR is merged
func (s *ApiSession) SetHubSessionKey(hubSessionKey string) {
	s.hubSessionKey = hubSessionKey
}
func (s *ApiSession) GetHubSessionKey() string {
	return s.hubSessionKey
}

func (s *ApiSession) SetUsernameAndPassword(username, password string) {
	s.username = username
	s.password = password
}

func (s *ApiSession) GetUsernameAndPassword() (string, string) {
	return s.username, s.password
}

func (s *ApiSession) AddServerURLforServerKey(serverURL, serverSessionKey string) {
	s.userServerURLByKey[serverSessionKey] = serverURL
}

func (s *ApiSession) GetServerURLbyServerKey(serverSessionKey string) string {
	return s.userServerURLByKey[serverSessionKey]
}

// New returns a new apiSession struct
func New() *ApiSession {

	return &ApiSession{
		hubSessionKey:      "",
		username:           "",
		password:           "",
		userServerURLByKey: make(map[string]string),
	}
}
