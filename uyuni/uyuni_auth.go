package uyuni

//Server authenticator
const (
	loginPath  = "auth.login"
	logoutPath = "auth.logout"
)

type uyuniAuthenticator struct {
	uyuniCallExecutor *uyuniCallExecutor
}

func NewUyuniAuthenticator(uyuniCallExecutor *uyuniCallExecutor) *uyuniAuthenticator {
	return &uyuniAuthenticator{uyuniCallExecutor}
}

func (u *uyuniAuthenticator) Login(endpoint, username, password string) (string, error) {
	response, err := u.uyuniCallExecutor.ExecuteCall(endpoint, loginPath, []interface{}{username, password})
	if err != nil {
		return "", err
	}
	return response.(string), nil
}

func (u *uyuniAuthenticator) Logout(endpoint, sessionKey string) error {
	_, err := u.uyuniCallExecutor.ExecuteCall(endpoint, logoutPath, []interface{}{sessionKey})
	if err != nil {
		return err
	}
	return nil
}
