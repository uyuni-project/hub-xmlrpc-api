package integration_tests

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/transformer"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
	"github.com/uyuni-project/hub-xmlrpc-api/uyuni"
)

func Test_Multicast(t *testing.T) {
	tt := []struct {
		name         string
		methodName   string
		call         string
		argsByServer map[int64][]interface{}
		output       string
	}{
		{name: "multicast.system.listSystems"},
		{
			name: "multicast.system.listUserSystems",
			argsByServer: map[int64][]interface{}{
				1: []interface{}{"admin", "admin"},
				2: []interface{}{"admin", "admin"},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.methodName, func(t *testing.T) {
			//setup env
			conf := config.NewConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.RequestTimeout)

			var syncMap sync.Map
			hubSessionRepository := session.NewInMemoryHubSessionRepository(&syncMap)
			serverSessionRepository := session.NewInMemoryServerSessionRepository(&syncMap)

			uyuniCallExecutor := uyuni.NewUyuniCallExecutor(client)
			uyuniAuthenticator := uyuni.NewUyuniAuthenticator(uyuniCallExecutor)
			uyuniTopologyInfoRetriever := uyuni.NewUyuniTopologyInfoRetriever(uyuniCallExecutor)

			serverAuthenticator := gateway.NewServerAuthenticator(conf.HubAPIURL, uyuniAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository, serverSessionRepository)
			hubTopologyInfoRetriever := gateway.NewTopologyInfoRetriever(conf.HubAPIURL, uyuniTopologyInfoRetriever)
			hubLoginer := gateway.NewHubLoginer(conf.HubAPIURL, uyuniAuthenticator, serverAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository)
			multicaster := gateway.NewMulticaster(uyuniCallExecutor, hubSessionRepository)

			loginController := controller.NewHubLoginController(hubLoginer, transformer.MulticastResponseTransformer)
			hubTopologyController := controller.NewHubTopologyController(hubTopologyInfoRetriever)
			multicastController := controller.NewMulticastController(multicaster, transformer.MulticastResponseTransformer)

			const xmlInput = `
			<methodCall>
			<methodName>%s</methodName>
			   <params>
				  <param>
					 <value><int>0000</int></value>
				  </param>
			   </params>
			</methodCall>`
			xmlBody := fmt.Sprintf(xmlInput, tc.methodName)

			//login
			req, err := http.NewRequest(http.MethodPost, conf.HubAPIURL, bytes.NewBuffer([]byte(xmlBody)))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			loginResponse := struct {
				Data *controller.LoginWithAutoconnectModeResponse
			}{}
			err = loginController.LoginWithAutoconnectMode(req, &controller.LoginRequest{"admin", "admin"}, &loginResponse)
			if err != nil {
				t.Fatalf("could not login to hub: %v", err)
			}

			//get the serverIDs
			hubSessionKey := loginResponse.Data.SessionKey
			listServerIDsRequest := struct{ HubSessionKey string }{hubSessionKey}
			listServerIDsResponse := struct{ Data []int64 }{}

			err = hubTopologyController.ListServerIDs(req, &listServerIDsRequest, &listServerIDsResponse)
			if err != nil {
				t.Fatalf("could not login to hub: %v", err)
			}

			//execute multicast call
			serverIDs := listServerIDsResponse.Data
			multicastRequest := controller.MulticastRequest{tc.call, hubSessionKey, serverIDs, tc.argsByServer}
			multicastReply := struct{ Data *controller.MulticastResponse }{}

			err = multicastController.Multicast(req, &multicastRequest, &multicastReply)

			if err != nil && tc.output != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
		})
	}
}
