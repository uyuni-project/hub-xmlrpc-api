package session

import (
	"reflect"
	"sync"
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
			var syncMap sync.Map
			repo := NewInMemoryHubSessionRepository(&syncMap)

			repo.SaveHubSession(tc.hubSession)
			hubSession := repo.RetrieveHubSession(tc.hubSessionKey)

			if !reflect.DeepEqual(hubSession, tc.hubSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.hubSession)
			}
		})
	}
}

func TestRetrieveHubSession(t *testing.T) {
	tt := []struct {
		name                   string
		hubSessionToSave       *gateway.HubSession
		hubSessionKeyToLookfor string
		expectedHubSession     *gateway.HubSession
	}{
		{name: "RetrieveHubSession Success",
			hubSessionToSave:       gateway.NewHubSession("sessionKey", "username", "password", 1),
			hubSessionKeyToLookfor: "sessionKey",
			expectedHubSession:     gateway.NewHubSession("sessionKey", "username", "password", 1),
		},
		{name: "RetrieveHubSession inexistent_hubSession_key",
			hubSessionToSave:       gateway.NewHubSession("sessionKey", "username", "password", 1),
			hubSessionKeyToLookfor: "inexistent_sessionKey",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var syncMap sync.Map
			repo := NewInMemoryHubSessionRepository(&syncMap)

			repo.SaveHubSession(tc.hubSessionToSave)

			hubSession := repo.RetrieveHubSession(tc.hubSessionKeyToLookfor)

			if !reflect.DeepEqual(hubSession, tc.expectedHubSession) {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", tc.expectedHubSession, hubSession)
			}
		})
	}
}

func TestSaveServerSessions(t *testing.T) {
	tt := []struct {
		name          string
		hubSessionKey string
		hubSession    *gateway.HubSession
		serverID      int64
		serverSession *gateway.ServerSession
	}{
		{name: "SaveServerSession Success",
			hubSessionKey: "sessionKey",
			hubSession:    gateway.NewHubSession("sessionKey", "username", "password", 1),
			serverID:      1234,
			serverSession: gateway.NewServerSession(1234, "url", "serverSessionKey", "sessionKey"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var syncMap sync.Map
			hubRepo := NewInMemoryHubSessionRepository(&syncMap)
			hubRepo.SaveHubSession(tc.hubSession)

			repo := NewInMemoryServerSessionRepository(&syncMap)

			expectedServerSessions := map[int64]*gateway.ServerSession{tc.serverID: tc.serverSession}

			repo.SaveServerSessions(tc.hubSessionKey, expectedServerSessions)

			serverSessions := repo.RetrieveServerSessions(tc.hubSessionKey)

			if !reflect.DeepEqual(serverSessions, expectedServerSessions) {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", expectedServerSessions, serverSessions)
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
			hubSessionToSave:       gateway.NewHubSession("sessionKey", "username", "password", 1),
			serverIDToSave:         1234,
			serverIDToLookfor:      1234,
			serverSessionToSave:    gateway.NewServerSession(1234, "url", "serverSessionKey", "sessionKey"),
			expectedServerSession:  gateway.NewServerSession(1234, "url", "serverSessionKey", "sessionKey"),
		},
		{name: "RetrieveServerSessionByServerID inexistent_hubSession_key",
			hubSessionKeyToSave:    "sessionKey",
			hubSessionKeyToLookfor: "inexistent_sessionKey",
			hubSessionToSave:       gateway.NewHubSession("sessionKey", "username", "password", 1),
			serverIDToSave:         1234,
			serverIDToLookfor:      1234,
			serverSessionToSave:    gateway.NewServerSession(1234, "url", "serverSessionKey", "sessionKey"),
			expectedServerSession:  nil,
		},
		{name: "RetrieveServerSessionByServerID inexistent_serverID",
			hubSessionKeyToSave:    "sessionKey",
			hubSessionKeyToLookfor: "sessionKey",
			hubSessionToSave:       gateway.NewHubSession("sessionKey", "username", "password", 1),
			serverIDToSave:         1234,
			serverIDToLookfor:      -1,
			serverSessionToSave:    gateway.NewServerSession(1234, "url", "serverSessionKey", "sessionKey"),
			expectedServerSession:  nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var syncMap sync.Map
			hubRepo := NewInMemoryHubSessionRepository(&syncMap)
			hubRepo.SaveHubSession(tc.hubSessionToSave)

			repo := NewInMemoryServerSessionRepository(&syncMap)

			repo.SaveServerSessions(tc.hubSessionKeyToSave, map[int64]*gateway.ServerSession{tc.serverIDToSave: tc.serverSessionToSave})

			serverSession := repo.RetrieveServerSessionByServerID(tc.hubSessionKeyToLookfor, tc.serverIDToLookfor)

			if !reflect.DeepEqual(serverSession, tc.expectedServerSession) {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", tc.expectedServerSession, serverSession)
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
			hubSession:            gateway.NewHubSession("sessionKey", "username", "password", 1),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var syncMap sync.Map
			repo := NewInMemoryHubSessionRepository(&syncMap)

			repo.SaveHubSession(tc.hubSession)
			hubSession := repo.RetrieveHubSession(tc.hubSessionKeyToSave)

			if !reflect.DeepEqual(hubSession, tc.hubSession) {
				t.Fatalf("expected and actual doesn't match, Expected was: %v", tc.hubSession)
			}

			repo.RemoveHubSession(tc.hubSessionKeyToRemove)

			hubSession = repo.RetrieveHubSession(tc.hubSessionKeyToRemove)

			if hubSession != nil {
				t.Fatalf("expected and actual doesn't match. Expected was:\n%v\nActual is:\n%v", nil, hubSession)
			}
		})
	}
}
