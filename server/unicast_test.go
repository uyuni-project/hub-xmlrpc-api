package server

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/uyuni-project/hub-xmlrpc-api/client"
	"github.com/uyuni-project/hub-xmlrpc-api/config"
	"github.com/uyuni-project/hub-xmlrpc-api/session"
)

func TestRemoveUnicastNamespace(t *testing.T) {

	tt := []struct {
		name   string
		input  string
		output string
	}{
		{name: "valid values-1", input: "unicast.list.servers", output: "list.servers"},
		{name: "valid values-2", input: "unicast.version", output: "version"},
		{name: "no-namespace", input: "version", output: ""},
		{name: "empty values", input: "", output: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			result := removeUnicastNamespace(tc.input)

			if result != tc.output {
				t.Fatalf("Unexpected result. Expected: %v, Got: %v", tc.output, result)
			}

		})
	}
}
func TestUniCastDefaultMethod(t *testing.T) {

	tt := []struct {
		name       string
		parameters []interface{}
		output     string
	}{
		{name: "unicast.system.listSystems"},
		{name: "unicast.system.listUserSystems", parameters: []interface{}{"admin"}},
		{name: "unicast.system.unknownmethod", parameters: []interface{}{"admin"}, output: "request error: bad status code - 400"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
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
			conf := config.InitializeConfig()
			hub := NewHubService(client.NewClient(conf), session.NewApiSession())
			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, bytes.NewBuffer([]byte(xmlBody)))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			credentials := struct {
				Username string
				Password string
			}{"admin", "admin"}
			reply := struct{ Data string }{""}
			//login
			hub.LoginWithAutoconnectMode(req, &credentials, &reply)
			sessionKey := struct{ HubSessionKey string }{reply.Data}
			//Get the server Ids
			serverIdsreply := struct{ Data []int64 }{}
			hub.ListServerIds(req, &sessionKey, &serverIdsreply)
			firstServerIDs := serverIdsreply.Data[0]
			unicastArgs := UnicastArgs{HubSessionKey: reply.Data, ServerID: firstServerIDs, ServerArgs: tc.parameters}
			unicastReply := struct{ Data interface{} }{}

			unicastService := NewUnicastService(hub.client, hub.apiSession)
			err = unicastService.DefaultMethod(req, &unicastArgs, &unicastReply)
			if err != nil {
				if tc.output != err.Error() {
					t.Fatalf("Error during executing request: %v", err)
				}
			}
		})
	}
}
