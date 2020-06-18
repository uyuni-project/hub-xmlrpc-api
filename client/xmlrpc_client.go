package client

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/kolo/xmlrpc"
)

type Client struct {
	connectTimeout, requestTimeout int
}

func NewClient(connectTimeout, requestTimeout int) *Client {
	return &Client{connectTimeout: connectTimeout, requestTimeout: requestTimeout}
}

func (c *Client) ExecuteCall(endpoint string, call string, args []interface{}) (response interface{}, err error) {
	client, err := getClientWithTimeout(endpoint, c.connectTimeout, c.requestTimeout)
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

func getClientWithTimeout(url string, connectTimeout, requestTimeout int) (*xmlrpc.Client, error) {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext:     timeoutDialer(time.Duration(connectTimeout)*time.Second, time.Duration(requestTimeout)*time.Second),
	}
	return xmlrpc.NewClient(url, &transport)
}
