package session

import (
	"reflect"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
)

func TestSaveHubSession(t *testing.T) {
	tt := []struct {
		name          string
		hubSessionKey string
		hubSession    *gateway.HubSession
	}{
		{name: "SaveHubSession Success",
			hubSessionKey: "sessionKey",
			hubSession:    gateway.NewHubSession("sessionKey", "username", "password", 1),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			session := NewHubSessionRepository()

			session.SaveHubSession(tc.hubSessionKey, tc.hubSession)
			hubSession := session.RetrieveHubSession(tc.hubSessionKey)

			if !reflect.DeepEqual(hubSession, tc.hubSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.hubSession)
			}
		})
	}
}

func TestRetrieveHubSession(t *testing.T) {
	tt := []struct {
		name                   string
		hubSessionKeyToSave    string
		hubSessionKeyToLookfor string
		hubSession             *gateway.HubSession
	}{
		{name: "RetrieveHubSession Success",
			hubSessionKeyToSave:    "sessionKey",
			hubSessionKeyToLookfor: "sessionKey",
			hubSession:             gateway.NewHubSession("username", "password", 1),
		},
		{name: "RetrieveHubSession inexistent_hubSession_key",
			hubSessionKeyToSave:    "sessionKey",
			hubSessionKeyToLookfor: "inexistent_sessionKey",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			session := NewSession()

			session.SaveHubSession(tc.hubSessionKeyToSave, tc.hubSession)

			hubSession := session.RetrieveHubSession(tc.hubSessionKeyToLookfor)

			if !reflect.DeepEqual(hubSession, tc.hubSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.hubSession)
			}
		})
	}
}

func TestSaveServerSession(t *testing.T) {
	tt := []struct {
		name          string
		hubSessionKey string
		hubSession    *gateway.HubSession
		serverID      int64
		serverSession *gateway.ServerSession
	}{
		{name: "SaveServerSession Success",
			hubSessionKey: "sessionKey",
			hubSession:    gateway.NewHubSession("username", "password", 1),
			serverID:      1234,
			serverSession: gateway.NewServerSession("url", "serverSessionKey"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			session := NewSession()

			session.SaveHubSession(tc.hubSessionKey, tc.hubSession)

			hubSession := session.RetrieveHubSession(tc.hubSessionKey)

			if !reflect.DeepEqual(hubSession, tc.hubSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.hubSessionKey)
			}

			session.SaveServerSession(tc.hubSessionKey, tc.serverID, tc.serverSession)

			serverSession := session.RetrieveServerSessionByServerID(tc.hubSessionKey, tc.serverID)

			if !reflect.DeepEqual(serverSession, tc.serverSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.serverSession)
			}
		})
	}
}

func TestRetrieveServerSessionByServerID(t *testing.T) {
	tt := []struct {
		name                   string
		hubSessionKeyToSave    string
		hubSessionKeyToLookfor string
		hubSessionToSave       *gateway.HubSession
		serverIDToSave         int64
		serverIDToLookfor      int64
		serverSessionToSave    *gateway.ServerSession
		expectedServerSession  *gateway.ServerSession
	}{
		{name: "RetrieveServerSessionByServerID Success",
			hubSessionKeyToSave:    "sessionKey",
			hubSessionKeyToLookfor: "sessionKey",
			hubSessionToSave:       gateway.NewHubSession("username", "password", 1),
			serverIDToSave:         1234,
			serverIDToLookfor:      1234,
			serverSessionToSave:    gateway.NewServerSession("url", "serverSessionKey"),
			expectedServerSession:  gateway.NewServerSession("url", "serverSessionKey"),
		},
		{name: "RetrieveServerSessionByServerID inexistent_hubSession_key",
			hubSessionKeyToSave:    "sessionKey",
			hubSessionKeyToLookfor: "inexistent_sessionKey",
			hubSessionToSave:       gateway.NewHubSession("username", "password", 1),
			serverIDToSave:         1234,
			serverIDToLookfor:      1234,
			serverSessionToSave:    gateway.NewServerSession("url", "serverSessionKey"),
			expectedServerSession:  nil,
		},
		{name: "RetrieveServerSessionByServerID inexistent_serverID",
			hubSessionKeyToSave:    "sessionKey",
			hubSessionKeyToLookfor: "sessionKey",
			hubSessionToSave:       gateway.NewHubSession("username", "password", 1),
			serverIDToSave:         1234,
			serverIDToLookfor:      -1,
			serverSessionToSave:    gateway.NewServerSession("url", "serverSessionKey"),
			expectedServerSession:  nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			session := NewSession()

			session.SaveHubSession(tc.hubSessionKeyToSave, tc.hubSessionToSave)
			session.SaveServerSession(tc.hubSessionKeyToSave, tc.serverIDToSave, tc.serverSessionToSave)

			serverSession := session.RetrieveServerSessionByServerID(tc.hubSessionKeyToLookfor, tc.serverIDToLookfor)

			if !reflect.DeepEqual(serverSession, tc.expectedServerSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.expectedServerSession)
			}
		})
	}
}

func TestRemoveHubSession(t *testing.T) {
	tt := []struct {
		name                  string
		hubSessionKeyToSave   string
		hubSessionKeyToRemove string
		hubSession            *gateway.HubSession
	}{
		{name: "RemoveHubSession Success",
			hubSessionKeyToSave:   "sessionKey",
			hubSessionKeyToRemove: "sessionKey",
			hubSession:            gateway.NewHubSession("username", "password", 1),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			session := NewSession()

			session.SaveHubSession(tc.hubSessionKeyToSave, tc.hubSession)

			hubSession := session.RetrieveHubSession(tc.hubSessionKeyToRemove)

			if !reflect.DeepEqual(hubSession, tc.hubSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.hubSession)
			}

			session.RemoveHubSession(tc.hubSessionKeyToRemove)

			hubSession = session.RetrieveHubSession(tc.hubSessionKeyToRemove)

			if hubSession != nil {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", nil)
			}
		})
	}
}
