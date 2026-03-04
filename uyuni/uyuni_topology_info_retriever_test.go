package uyuni

import (
	"errors"
	"reflect"
	"testing"
)

type mockClient struct {
	mockExecuteCall func(endpoint string, call string, args []interface{}) (response interface{}, err error)
	mockAddRootCA   func(pem []byte) error
}

func (m *mockClient) ExecuteCall(endpoint string, call string, args []interface{}) (interface{}, error) {
	return m.mockExecuteCall(endpoint, call, args)
}

func (m *mockClient) AddRootCA(pem []byte) error {
	if m.mockAddRootCA != nil {
		return m.mockAddRootCA(pem)
	}
	return nil
}

func TestListServerIDs(t *testing.T) {
	client := &mockClient{}
	executor := NewUyuniCallExecutor(client)
	retriever := NewUyuniTopologyInfoRetriever(executor, true)

	sessionKey := "test-session"
	endpoint := "http://hub/rpc/api"

	// Mock sync.hub.listPeripheralServers
	client.mockExecuteCall = func(ep string, call string, args []interface{}) (interface{}, error) {
		if call == "sync.hub.listPeripheralServers" {
			return []interface{}{
				map[string]interface{}{
					"id":      int64(1001),
					"fqdn":    "peripheral1.example.com",
					"root_ca": "PEM1",
				},
				map[string]interface{}{
					"id":      int64(1002),
					"fqdn":    "peripheral2.example.com",
					"root_ca": "PEM2",
				},
			}, nil
		}
		return nil, nil
	}

	rootCAsAdded := make([]string, 0)
	client.mockAddRootCA = func(pem []byte) error {
		rootCAsAdded = append(rootCAsAdded, string(pem))
		return nil
	}

	ids, err := retriever.ListServerIDs(endpoint, sessionKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []int64{1001, 1002}
	if !reflect.DeepEqual(ids, expectedIDs) {
		t.Errorf("expected IDs %v, got %v", expectedIDs, ids)
	}

	// Verify cache
	retriever.mu.RLock()
	endpoint1, ok1 := retriever.serverEndpoints[1001]
	endpoint2, ok2 := retriever.serverEndpoints[1002]
	retriever.mu.RUnlock()

	if !ok1 || endpoint1 != "https://peripheral1.example.com/rpc/api" {
		t.Errorf("unexpected endpoint for 1001: %v", endpoint1)
	}
	if !ok2 || endpoint2 != "https://peripheral2.example.com/rpc/api" {
		t.Errorf("unexpected endpoint for 1002: %v", endpoint2)
	}

	// Verify root CAs
	expectedRootCAs := []string{"PEM1", "PEM2"}
	if !reflect.DeepEqual(rootCAsAdded, expectedRootCAs) {
		t.Errorf("expected root CAs %v, got %v", expectedRootCAs, rootCAsAdded)
	}
}

func TestListServerIDsFallback(t *testing.T) {
	client := &mockClient{}
	executor := NewUyuniCallExecutor(client)
	retriever := NewUyuniTopologyInfoRetriever(executor, true)

	sessionKey := "test-session"
	endpoint := "http://hub/rpc/api"

	// Mock sync.hub.listPeripheralServers to fail
	client.mockExecuteCall = func(ep string, call string, args []interface{}) (interface{}, error) {
		if call == "sync.hub.listPeripheralServers" {
			return nil, errors.New("fail")
		}
		if call == listSystemsWithEntitlementPath {
			return []interface{}{
				map[string]interface{}{"id": int64(2001)},
			}, nil
		}
		return nil, nil
	}

	ids, err := retriever.ListServerIDs(endpoint, sessionKey)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedIDs := []int64{2001}
	if !reflect.DeepEqual(ids, expectedIDs) {
		t.Errorf("expected IDs %v, got %v", expectedIDs, ids)
	}
}
