package integration_tests

// import (
// 	"os"
// 	"regexp"
// 	"strings"
// 	"sync"
// 	"testing"

// 	"github.com/uyuni-project/hub-xmlrpc-api/client"
// 	"github.com/uyuni-project/hub-xmlrpc-api/config"
// 	"github.com/uyuni-project/hub-xmlrpc-api/controller"
// 	"github.com/uyuni-project/hub-xmlrpc-api/controller/transformer"
// 	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
// 	"github.com/uyuni-project/hub-xmlrpc-api/session"
// 	"github.com/uyuni-project/hub-xmlrpc-api/uyuni"
// )

// func init() {
// 	// load test data
// 	os.Setenv("HUB_CONFIG_FILE", "../tests/config.json")
// }

// func TestLogin(t *testing.T) {
// 	tt := []struct {
// 		name     string
// 		username string
// 		password string
// 		err      string
// 	}{
// 		{name: "Invalid credentials", username: "unknown-user", password: "unknown-user", err: controller.FaultInvalidCredentials.Message},
// 		{name: "Valid credentials", username: "admin", password: "admin"},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			//setup env
// 			conf := config.NewConfig()
// 			client := client.NewClient(conf.ConnectTimeout, conf.RequestTimeout)

// 			var syncMap sync.Map
// 			hubSessionRepository := session.NewInMemoryHubSessionRepository(&syncMap)
// 			serverSessionRepository := session.NewInMemoryServerSessionRepository(&syncMap)

// 			uyuniCallExecutor := uyuni.NewUyuniCallExecutor(client)
// 			uyuniAuthenticator := uyuni.NewUyuniAuthenticator(uyuniCallExecutor)
// 			uyuniTopologyInfoRetriever := uyuni.NewUyuniTopologyInfoRetriever(uyuniCallExecutor)

// 			serverAuthenticator := gateway.NewServerAuthenticator(conf.HubAPIURL, uyuniAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository, serverSessionRepository)
// 			hubLoginer := gateway.NewHubLoginer(conf.HubAPIURL, uyuniAuthenticator, serverAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository)

// 			loginController := controller.NewHubLoginController(hubLoginer, transformer.MulticastResponseTransformer)

// 			reply := struct{ Data string }{""}
// 			err := loginController.Login(nil, &controller.LoginRequest{tc.username, tc.password}, &reply)

// 			if err != nil && !strings.Contains(err.Error(), tc.err) {
// 				t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
// 			}

// 			// test the hubkey
// 			hubsessionkey := reply.Data

// 			matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, hubsessionkey)
// 			if !matched {
// 				t.Fatalf("Unexepected token pattern %v", hubsessionkey)
// 				return
// 			}

// 			hubSession := hubSessionRepository.RetrieveHubSession(hubsessionkey)
// 			if hubSession.HubSessionKey != hubsessionkey {
// 				t.Fatalf("User name doesn't match with the key, expected %v, got %v", hubSession.HubSessionKey, hubsessionkey)
// 			}
// 		})
// 	}

// }

// func TestLoginAutoconnect(t *testing.T) {
// 	tt := []struct {
// 		name     string
// 		username string
// 		password string
// 		err      string
// 	}{
// 		{name: "Valid credentials", username: "admin", password: "admin"},
// 		{name: "Invalid credentials", username: "unknown-user", password: "unknown-user", err: controller.FaultInvalidCredentials.Message},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			//setup env
// 			conf := config.NewConfig()
// 			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)

// 			var syncMap sync.Map
// 			hubSessionRepository := session.NewInMemoryHubSessionRepository(&syncMap)
// 			serverSessionRepository := session.NewInMemoryServerSessionRepository(&syncMap)

// 			uyuniCallExecutor := uyuni.NewUyuniCallExecutor(client)
// 			uyuniAuthenticator := uyuni.NewUyuniAuthenticator(uyuniCallExecutor)
// 			uyuniTopologyInfoRetriever := uyuni.NewUyuniTopologyInfoRetriever(uyuniCallExecutor)

// 			serverAuthenticator := gateway.NewServerAuthenticator(conf.Hub.SUMA_API_URL, uyuniAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository, serverSessionRepository)
// 			hubLoginer := gateway.NewHubLoginer(conf.Hub.SUMA_API_URL, uyuniAuthenticator, serverAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository)

// 			loginController := controller.NewHubLoginController(hubLoginer, transformer.MulticastResponseTransformer)

// 			reply := struct {
// 				Data *controller.LoginWithAutoconnectModeResponse
// 			}{}
// 			err := loginController.LoginWithAutoconnectMode(nil, &controller.LoginRequest{tc.username, tc.password}, &reply)
// 			if err != nil && !strings.Contains(err.Error(), tc.err) {
// 				t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
// 			}

// 			// test the hubkey
// 			hubSessionKey := reply.Data.SessionKey

// 			matched, _ := regexp.MatchString(`^[A-Za-z0-9]{68}$`, hubSessionKey)
// 			if !matched {
// 				t.Fatalf("Unexepected token pattern %v", hubSessionKey)
// 				return
// 			}

// 			hubSession := hubSessionRepository.RetrieveHubSession(hubSessionKey)
// 			if hubSession.HubSessionKey != hubSessionKey {
// 				t.Fatalf("User name doesn't match with the key, expected %v, got %v", hubSession.HubSessionKey, hubSessionKey)
// 			}

// 			//test if servers attached to hub have also been authenticated automatically
// 			serverIDs := struct{ Data []int64 }{}
// 			err = uyuniTopologyInfoRetriever.ListServerIDs(hubSessionKey), &serverIDs)

// 			for _, serverID := range serverIDs.Data {
// 				serverSession := serverSessionRepository.RetrieveServerSessionByServerID(hubSessionKey, serverID)
// 				if len(serverSession.url) == 0 {
// 					t.Fatalf("Expected valid url for server with severId: %v, got empty instead %v", serverID, serverSession.url)
// 				}
// 				if len(serverSession.sessionkey) <= 0 {
// 					t.Fatalf("Expected valid SessionKey for server with severId: %v, Got %v", serverID, serverSession.sessionkey)
// 				}
// 			}
// 		})
// 	}
// }

// func TestAttachToServers(t *testing.T) {
// 	tt := []struct {
// 		name     string
// 		username string
// 		password string
// 		err      string
// 	}{
// 		{name: "Valid credentials", username: "admin", password: "admin"},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			conf := config.InitializeConfig()
// 			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
// 			session := session.NewSession()
// 			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

// 			//login
// 			reply := struct{ Data string }{""}
// 			err = hub.LoginWithAuthRelayMode(nil, &server.LoginArgs{tc.username, tc.password}, &reply)
// 			if err != nil {
// 				t.Fatalf("Login faied with error : %v", err)
// 			}
// 			sessionKey := struct{ HubSessionKey string }{reply.Data}

// 			// List server Ids
// 			serverIDsReply := struct{ Data []int64 }{}
// 			err = hub.ListServerIds(req, &sessionKey, &serverIDsReply)
// 			serverIDs := serverIDsReply.Data

// 			srvArgs := server.MulticastArgs{"method", sessionKey.HubSessionKey, serverIDs, nil}
// 			err = hub.AttachToServers(req, &srvArgs, &struct{ Data []error }{})
// 			if err != nil && err.Error() != tc.err {
// 				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.err, err.Error())
// 				return
// 			}
// 			for _, s := range serverIDs {
// 				serverSession := session.RetrieveServerSessionByServerID(sessionKey.HubSessionKey, s)
// 				if len(serverSession.url) == 0 {
// 					t.Fatalf("Expected valid url for server with severId: %v, got empty instead %v", s, serverSession.url)
// 				}
// 				if len(serverSession.sessionKey) <= 0 {
// 					t.Fatalf("Expected valid SessionKey for server with severId: %v, Got %v", s, serverSession.sessionKey)
// 				}
// 			}
// 		})
// 	}
// }

// func TestListServerIds(t *testing.T) {
// 	tt := []struct {
// 		name     string
// 		username string
// 		password string
// 		err      string
// 	}{
// 		{name: "Valid credentials", username: "admin", password: "admin"},
// 		{name: "With invalid  credentials", username: "unknownadmin", password: "unknownadmin", err: codec.FaultInvalidCredentials.Message},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			conf := config.InitializeConfig()
// 			client := client.NewClient(conf.ConnectTimeout, conf.ReadWriteTimeout)
// 			session := session.NewSession()
// 			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)

// 			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
// 			if err != nil {
// 				t.Fatalf("could not create request: %v", err)
// 			}
// 			reply := struct{ Data string }{""}
// 			err = hub.Login(req, &server.LoginArgs{tc.username, tc.password}, &reply)
// 			if err != nil {
// 				if !strings.Contains(err.Error(), tc.err) {
// 					t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
// 				}
// 				return
// 			}
// 			serverIdsreply := struct{ Data []int64 }{}
// 			err = hub.ListServerIds(req, &struct{ HubSessionKey string }{reply.Data}, &serverIdsreply)
// 			if err != nil && err.Error() != tc.err {
// 				t.Fatalf("Unexpected Result: Exepected %v, Got %v", tc.err, err.Error())
// 				return
// 			}
// 			serverIds := len(serverIdsreply.Data)
// 			if serverIds <= 0 {
// 				t.Fatalf("Unexpected Result: Expected some servers, got nothing")
// 			}
// 		})
// 	}
// }
