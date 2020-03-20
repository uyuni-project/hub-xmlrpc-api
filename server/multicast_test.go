package server

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestGetKeysAndValuesFromMap(t *testing.T) {
	input1 := make(map[int64]interface{})
	input1[1000010000] = "value1"
	input1[1000010001] = "value2"
	expectedKeysForInput1 := []int64{1000010000, 1000010001}
	expectedValuesForInput1 := []interface{}{"value1", "value2"}
	input2 := make(map[int64]interface{})
	input2[1000010000] = "value1"
	expectedKeysForInput2 := []int64{1000010000}
	expectedValuesForInput2 := []interface{}{"value1"}
	input3 := make(map[int64]interface{})
	expectedKeysForInput3 := []int64{}
	expectedValuesForInput3 := []interface{}{}

	tt := []struct {
		name           string
		values         map[int64]interface{}
		expectedKeys   []int64
		expectedValues []interface{}
	}{
		{name: "valid values-1", values: input1, expectedKeys: expectedKeysForInput1, expectedValues: expectedValuesForInput1},
		{name: "valid values-2", values: input2, expectedKeys: expectedKeysForInput2, expectedValues: expectedValuesForInput2},
		{name: "empty values", values: input3, expectedKeys: expectedKeysForInput3, expectedValues: expectedValuesForInput3},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			successfulKeys, successfulValues := getKeysAndValuesFromMap(tc.values)

			if !reflect.DeepEqual(successfulKeys, tc.expectedKeys) {
				t.Fatalf("Not equal")
			}
			if !reflect.DeepEqual(successfulValues, tc.expectedValues) {
				t.Fatalf("Not equal")
			}
		})
	}
}
func TestRemoveMulticastNamespace(t *testing.T) {

	tt := []struct {
		name   string
		input  string
		output string
	}{
		{name: "valid values-1", input: "multicast.list.servers", output: "list.servers"},
		{name: "valid values-2", input: "multicast.version", output: "version"},
		{name: "empty values", input: "", output: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			result := removeMulticastNamespace(tc.input)

			if result != tc.output {
				t.Fatalf("Unexpected result. Expected: %v, Got: %v", tc.output, result)
			}

		})
	}
}

func TestResolveMulticastServerArgs(t *testing.T) {
	input := [][]interface{}{[]interface{}{"param1-server1", "param1-server2"}, []interface{}{"param2-server1", "param2-server2"}}
	srvArgoutput := [][]interface{}{{"param1-server1", "param2-server1"}, {"param1-server2", "param2-server2"}}
	tt := []struct {
		name     string
		username string
		password string
		data     [][]interface{}
		output   [][]interface{}
		err      string
	}{
		{name: "valid-values", username: "admin", password: "admin", data: input, output: srvArgoutput},
		{name: "empty-values", username: "unknownuser", password: "unknownuser", err: FaultInvalidCredentials.String},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hub := Hub{}
			req, err := http.NewRequest("GET", conf.Hub.SUMA_API_URL, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			reply := struct{ Data string }{""}
			err = new(Hub).Login(req, &struct{ Username, Password string }{tc.username, tc.password}, &reply)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Expected %v, Got %v", tc.err, err.Error())
				}
				return
			}
			sessionKey := struct{ HubSessionKey string }{reply.Data}
			serverIdsreply := struct{ Data []int64 }{}
			err = hub.ListServerIds(req, &sessionKey, &serverIdsreply)
			if err != nil {
				t.Fatalf("could not get the sever ids : %v", err)
			}
			serverIds := serverIdsreply.Data

			srvArgs := MulticastArgs{sessionKey.HubSessionKey, serverIds, tc.data}
			result := resolveMulticastServerArgs(&srvArgs)
			resultLength := len(result)
			if resultLength != len(serverIds) {
				t.Fatalf("Unexpected result. Length should be same but got %v & %v", resultLength, len(serverIds))
			}
			serverIdsFromResult := make([]int64, resultLength)
			for i, v := range result {
				serverIdsFromResult[i] = v.serverID
			}
			if !reflect.DeepEqual(serverIdsFromResult, serverIds) {
				t.Fatalf("Unexpected result: server Ids are not same")
			}
			for i, v := range tc.output {
				//get the matching result
				res := result[i]
				// remove the 1st element from the array and then compare as first is sessionkey
				if !reflect.DeepEqual(v, res.args[1:len(res.args)]) {
					t.Fatalf("Unexpected result: expected %v, got %v", tc.output, res.args)
				}
			}

		})
	}
}

func TestMulticastCall(t *testing.T) {
	input := [][]interface{}{[]interface{}{"admin", "admin"}}
	//srvArgoutput := [][]interface{}{{"param1-server1", "param2-server1"}, {"param1-server2", "param2-server2"}}
	tt := []struct {
		name       string
		parameters [][]interface{}
		output     string
	}{
		{name: "system.listSystems"},
		{name: "system.listUserSystems", parameters: input},
		{name: "system.unknownmethod", parameters: input, output: "request error: bad status code - 400"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hub := Hub{}
			req, err := http.NewRequest("GET", "localhost:8888", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			//Login
			reply := struct{ Data string }{}
			err = hub.LoginWithAutoconnectMode(req, &struct{ Username, Password string }{"admin", "admin"}, &reply)
			if err != nil {
				t.Fatalf("Login fails: %v", err)
			}
			sessionKey := struct{ HubSessionKey string }{reply.Data}
			//Get the server Ids
			serverIdsreply := struct{ Data []int64 }{}
			hub.ListServerIds(req, &sessionKey, &serverIdsreply)
			serverIds := serverIdsreply.Data

			srvArgs := MulticastArgs{sessionKey.HubSessionKey, serverIds, tc.parameters}

			result := resolveMulticastServerArgs(&srvArgs)

			multicastResponse := multicastCall(tc.name, result)

			failedResponses := len(multicastResponse.Failed.Responses)
			successfulResponses := len(multicastResponse.Successfull.Responses)
			totalResponses := failedResponses + successfulResponses
			if totalResponses != len(serverIds) {
				t.Fatalf("Results are not complete, there should be result for every server. Expected number of reponse %v , Got %v", len(serverIds), totalResponses)
			}
			//t.Fatalf("sss %v", multicastResponse)
			if tc.name == "system.unknownmethod" {
				if failedResponses != len(serverIds) {
					t.Fatalf("Expected all response to come as failed but that didn't happen")
				}
				if multicastResponse.Failed.Responses[0] != tc.output {
					t.Fatalf("Expected %v, Got %v", multicastResponse.Failed.Responses[0], tc.output)

				}

			}
		})
	}
}
