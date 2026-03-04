package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	xmlrpc "github.com/uyuni-project/xmlrpc-public-methods"
)

type Client struct {
	connectTimeout, requestTimeout int
	rootCAs                        *x509.CertPool
	mu                             sync.RWMutex
}

func NewClient(connectTimeout, requestTimeout int) *Client {
	return &Client{
		connectTimeout: connectTimeout,
		requestTimeout: requestTimeout,
		rootCAs:        x509.NewCertPool(),
	}
}

func (c *Client) AddRootCA(pem []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.rootCAs.AppendCertsFromPEM(pem) {
		return fmt.Errorf("failed to parse root CA certificate")
	}
	return nil
}

func (c *Client) ExecuteCall(endpoint string, call string, args []interface{}) (response interface{}, err error) {
	client, err := c.getClientWithTimeout(endpoint, c.connectTimeout, c.requestTimeout)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	err = client.Call(call, args, &response)
	return response, err
}

func timeoutDialer(connectTimeout, requestTimeout time.Duration) func(ctx context.Context, net, addr string) (c net.Conn, err error) {
	return func(ctx context.Context, netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(requestTimeout))
		return conn, nil
	}
}

func (c *Client) getClientWithTimeout(url string, connectTimeout, requestTimeout int) (*xmlrpc.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	transport := http.Transport{
		DialContext: timeoutDialer(time.Duration(connectTimeout)*time.Second, time.Duration(requestTimeout)*time.Second),
		TLSClientConfig: &tls.Config{
			RootCAs: c.rootCAs,
		},
	}
	return xmlrpc.NewClient(url, &transport)
}
