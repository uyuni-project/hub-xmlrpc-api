package integration_tests

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/controller"
	"github.com/uyuni-project/hub-xmlrpc-api/controller/transformer"
	"github.com/uyuni-project/hub-xmlrpc-api/gateway"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
	"github.com/uyuni-project/hub-xmlrpc-api/uyuni"
)

func TestUniCast(t *testing.T) {
	tt := []struct {
		name             string
		methodName       string
		loginCredentials struct{ username, password string }
		args             []interface{}
		output           string
	}{
		{
			methodName:       "unicast.system.listSystems",
			loginCredentials: struct{ username, password string }{"admin", "admin"},
			args:             []interface{}{},
		},
		// {methodName: "unicast.system.listUserSystems", args: []interface{}{"admin"}},
		// {methodName: "unicast.system.unknownmethod", args: []interface{}{"admin"}, output: "request error: bad status code - 400"},
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
			hubLoginer := gateway.NewHubLoginer(conf.HubAPIURL, uyuniAuthenticator, serverAuthenticator, uyuniTopologyInfoRetriever, hubSessionRepository)
			unicaster := gateway.NewUnicaster(uyuniCallExecutor, serverSessionRepository)

			loginController := controller.NewHubLoginController(hubLoginer, transformer.MulticastResponseTransformer)
			unicastController := controller.NewUnicastController(unicaster)

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
			req, err := http.NewRequest("POST", conf.HubAPIURL, bytes.NewBuffer([]byte(xmlBody)))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			loginResponse := struct {
				Data *controller.LoginWithAutoconnectModeResponse
			}{}
			loginRequest := &controller.LoginRequest{tc.loginCredentials.username, tc.loginCredentials.password}
			err = loginController.LoginWithAutoconnectMode(req, loginRequest, &loginResponse)
			if err != nil {
				t.Fatalf("could not login to hub: %v", err)
			}

			serverID := loginResponse.Data.Successful.ServerIds[0]
			//execute unicast call
			unicastRequest := controller.UnicastRequest{loginResponse.Data.SessionKey, tc.methodName, serverID, tc.args}
			unicastResponse := struct{ Data interface{} }{}

			err = unicastController.Unicast(req, &unicastRequest, &unicastResponse)

			if err != nil && tc.output != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
		})
	}
}
