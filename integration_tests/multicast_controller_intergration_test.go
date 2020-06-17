package integration_tests

/*
import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/server"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
)
func Test_MulticastService_DefaultMethod(t *testing.T) {
	tt := []struct {
		name       string
		parameters [][]interface{}
		output     string
	}{
		{name: "multicast.system.listSystems"},
		{name: "multicast.system.listUserSystems", parameters: [][]interface{}{{"admin", "admin"}}},
		{name: "multicast.system.unknownmethod", parameters: [][]interface{}{{"admin", "admin"}}, output: "request error: bad status code - 400"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//setup env
			conf := config.InitializeConfig()
			client := client.NewClient(conf.ConnectTimeout, conf.RequestTimeout)
			session := session.NewSession()
			hub := server.NewHubService(client, session, conf.Hub.SUMA_API_URL)
			const xmlInput = `
			<methodCall>
			<methodName>%s</methodName>
			   <params>
				  <param>
					 <value><int>0000</int></value>
				  </param>
			   </params>
			</methodCall>`
			xmlBody := fmt.Sprintf(xmlInput, tc.name)
			//login
			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, bytes.NewBuffer([]byte(xmlBody)))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			loginReply := struct{ Data string }{""}
			err = hub.LoginWithAutoconnectMode(req, &server.LoginArgs{"admin", "admin"}, &loginReply)
			if err != nil {
				t.Fatalf("could not login to hub: %v", err)
			}
			//get the serverIDs
			hubsessionKey := struct{ HubSessionKey string }{loginReply.Data}
			serverIDs := struct{ Data []int64 }{}
			err = hub.ListServerIds(req, &hubsessionKey, &serverIDs)
			if err != nil {
				t.Fatalf("could not login to hub: %v", err)
			}
			//execute multicast call
			multicastArgs := server.MulticastArgs{tc.name, hubsessionKey.HubSessionKey, serverIDs.Data, tc.parameters}
			multicastReply := struct{ Data server.MulticastResponse }{}
			multicastService := server.NewMulticastService(client, session, conf.Hub.SUMA_API_URL)
			err = multicastService.DefaultMethod(req, &multicastArgs, &multicastReply)
			if err != nil && tc.output != err.Error() {
				t.Fatalf("Error during executing request: %v", err)
			}
		})
	}
}
*/
