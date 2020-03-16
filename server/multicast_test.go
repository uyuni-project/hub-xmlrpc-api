package server

import (
	"net/http"
	"reflect"
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
		name       string
		sessionkey string
		data       [][]interface{}
		output     [][]interface{}
	}{
		{name: "valid-values", sessionkey: SESSIONKEY, data: input, output: srvArgoutput},
		{name: "empty-values", sessionkey: SESSIONKEY},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hub := Hub{}
			req, err := http.NewRequest("GET", "localhost:8888", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			sessionKey := struct{ HubSessionKey string }{tc.sessionkey}
			serverIdsreply := struct{ Data []int64 }{}
			hub.ListServerIds(req, &sessionKey, &serverIdsreply)
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
