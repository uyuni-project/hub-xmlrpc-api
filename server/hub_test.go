package server

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

const SESSIONKEY = "300x2413800c14c02928568674dad9e71e0f061e2920be1d7c6542683d78de524bd4"

func init() {
	/* load test data */
	os.Setenv("HUB_CONFIG_FILE", "../tests/config.json")
	InitConfig()
}
func TestIsHubSessionValid(t *testing.T) {
	const errorMessage = "is not valid"
	tt := []struct {
		name       string
		sessionkey string
		result     bool
	}{
		{name: "Valid Session Key", sessionkey: SESSIONKEY, result: true},
		{name: "Invalid  Session Key", sessionkey: "300x241bad-key", result: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			isvalid := isHubSessionValid(tc.sessionkey)
			if isvalid != tc.result {
				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.result, isvalid)
			}
		})
	}
}

func TestListServerIds(t *testing.T) {

	const errorMessage = "Provided session key is invalid."
	tt := []struct {
		name       string
		sessionkey string
		err        string
	}{
		{name: "With valid Session Key", sessionkey: SESSIONKEY, err: ""},
		{name: "With invalid  Session Key", sessionkey: "300x241bad-key", err: errorMessage},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8888", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			hub := Hub{}
			sessionKey := struct{ HubSessionKey string }{tc.sessionkey}
			reply := struct{ Data []int64 }{}

			err = hub.ListServerIds(req, &sessionKey, &reply)
			if err != nil && err.Error() != tc.err {
				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.err, err.Error())
				return
			}
			serverIds := len(reply.Data)
			t.Logf("Number of returned servers: %v", serverIds)

		})
	}

}

func TestAreAllArgumentsOfSameLength(t *testing.T) {
	//ToDo: Make it better
	var sub1 = []string{"a", "b", "c", "d"}
	var sub2 = []string{"a", "b"}
	var main = [][]string{sub1, sub2}
	s := make([]interface{}, len(sub1))
	for i, v := range sub1 {
		s[i] = v
	}
	s1 := make([]interface{}, len(sub2))
	for i, v := range sub2 {
		s1[i] = v
	}
	t1 := make([][]interface{}, len(main))
	t1[0] = s
	t1[1] = s1

	fmt.Println(areAllArgumentsOfSameLength(t1))
	if areAllArgumentsOfSameLength(t1) != false {
		t.Fatalf("expected and actual doesn't match, Expected was: %v", false)
	}

	sub1 = []string{"a", "b", "c", "d"}
	sub2 = []string{"a", "b", "c", "e"}
	main = [][]string{sub1, sub2}
	s = make([]interface{}, len(sub1))
	for i, v := range sub1 {
		s[i] = v
	}
	s1 = make([]interface{}, len(sub2))
	for i, v := range sub2 {
		s1[i] = v
	}
	t1 = make([][]interface{}, len(main))
	t1[0] = s
	t1[1] = s1

	fmt.Println(areAllArgumentsOfSameLength(t1))
	if areAllArgumentsOfSameLength(t1) != true {
		t.Fatalf("expected and actual doesn't match, Expected was: %v", true)
	}
}

func TestLogin(t *testing.T) {

	//config.InitializeConfig()
	const errorMessage = "Either the password or username is incorrect."
	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Invalid credentials", username: "falsadmin", password: "falsadmin", err: errorMessage},
		{name: "Valid credentials", username: "admin", password: "admin"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8888", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			hub := Hub{}
			credentials := struct {
				Username string
				Password string
			}{tc.username, tc.password}
			reply := struct{ Data string }{""}

			err = hub.Login(req, &credentials, &reply)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Error message `%v` doesn't contain `%v`", err, tc.err)
					return
				}
			} else {

				matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, reply.Data)
				if !matched {
					t.Fatalf("Unexepected token pattern %v", reply.Data)
					return
				}
				username, password := apiSession.GetUsernameAndPassword(reply.Data)
				if username != tc.username {
					t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.username, username)
				}
				if password != tc.password {
					t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.password, password)
				}

			}
		})
	}

}
func TestLoginAutoconnect(t *testing.T) {
	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8888", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			hub := Hub{}
			credentials := struct {
				Username string
				Password string
			}{tc.username, tc.password}
			reply := struct{ Data string }{""}

			err = hub.LoginWithAutoconnectMode(req, &credentials, &reply)
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("Expected %v, Got %v", tc.err, err)
					return
				}
			} else {

				matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, reply.Data)
				if !matched {
					t.Fatalf("Unexepected token pattern %v", reply.Data)
					return
				}
				username, password := apiSession.GetUsernameAndPassword(reply.Data)
				if username != tc.username {
					t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.username, username)
				}
				if password != tc.password {
					t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.password, password)
				}

			}
		})
	}

}
func TestLoginWithAuthRelayMode(t *testing.T) {
	tt := []struct {
		name     string
		username string
		password string
		err      string
	}{
		{name: "Valid credentials", username: "admin", password: "admin"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8888", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			hub := Hub{}
			credentials := struct {
				Username string
				Password string
			}{tc.username, tc.password}
			reply := struct{ Data string }{""}

			err = hub.LoginWithAutoconnectMode(req, &credentials, &reply)
			if err != nil {
				if err.Error() != tc.err {
					t.Fatalf("Expected %v, Got %v", tc.err, err)
					return
				}
			} else {

				matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, reply.Data)
				if !matched {
					t.Fatalf("Unexepected token pattern %v", reply.Data)
					return
				}
				username, password := apiSession.GetUsernameAndPassword(reply.Data)
				if username != tc.username {
					t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.username, username)
				}
				if password != tc.password {
					t.Fatalf("User name doesn't match with the key, expected %v, got %v", tc.password, password)
				}

			}
		})
	}

}

func TestAttachToServers(t *testing.T) {

	const errorMessage = "Provided session key is invalid."
	tt := []struct {
		name       string
		sessionkey string
		err        string
	}{
		{name: "With valid Session Key", sessionkey: SESSIONKEY, err: ""},
		{name: "With invalid  Session Key", sessionkey: "300x241bad-key", err: errorMessage},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8888", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			hub := Hub{}
			sessionKey := struct{ HubSessionKey string }{tc.sessionkey}
			reply := struct{ Data []error }{}
			usernames := []interface{}{"admin", "admin"}
			passwords := []interface{}{"admin", "admin"}

			serverIdsreply := struct{ Data []int64 }{}
			err = hub.ListServerIds(req, &sessionKey, &serverIdsreply)
			serverIds := serverIdsreply.Data

			srvArgs := MulticastArgs{sessionKey.HubSessionKey, serverIds, [][]interface{}{usernames, passwords}}
			err = hub.AttachToServers(req, &srvArgs, &reply)
			if err != nil && err.Error() != tc.err {
				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.err, err.Error())
				return
			}
			for _, s := range serverIds {
				url, sesskey := apiSession.GetServerSessionInfoByServerID(sessionKey.HubSessionKey, s)
				if len(url) == 0 {
					t.Fatalf("Expected valid url for server with severId: %v, got empty instead %v", s, url)
				}
				if len(sesskey) <= 0 {
					t.Fatalf("Expected valid SessionKey for server with severId: %v, Got %v", s, sesskey)
				}

			}

			t.Logf("Number of returned servers: %v", serverIds)

		})
	}

}
