package uyuni

//Server authenticator
const (
	loginPath  = "auth.login"
	logoutPath = "auth.logout"
)

type uyuniServerAuthenticator struct {
	uyuniServerCallExecutor *uyuniServerCallExecutor
}

func NewUyuniServerAuthenticator(uyuniServerCallExecutor *uyuniServerCallExecutor) *uyuniServerAuthenticator {
	return &uyuniServerAuthenticator{uyuniServerCallExecutor}
}

func (u *uyuniServerAuthenticator) Login(endpoint, username, password string) (string, error) {
	response, err := u.uyuniServerCallExecutor.ExecuteCall(endpoint, loginPath, []interface{}{username, password})
	if err != nil {
		return "", err
	}
	return response.(string), nil
}

func (u *uyuniServerAuthenticator) Logout(endpoint, sessionKey string) error {
	_, err := u.uyuniServerCallExecutor.ExecuteCall(endpoint, logoutPath, []interface{}{sessionKey})
	if err != nil {
		return err
	}
	return nil
}

//Server call executor
type uyuniServerCallExecutor struct {
	client Client
}

func NewUyuniServerCallExecutor(client Client) *uyuniServerCallExecutor {
	return &uyuniServerCallExecutor{client}
}

func (u *uyuniServerCallExecutor) ExecuteCall(endpoint, call string, args []interface{}) (interface{}, error) {
	response, err := u.client.ExecuteCall(endpoint, call, args)
	if err != nil {
		return "", err
	}
	return response, nil
}
