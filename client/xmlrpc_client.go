package client

import (
	"net"
	"net/http"
	"time"

	xmlrpc "github.com/uyuni-project/xmlrpc-public-methods"
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

func timeoutDialer(connectTimeout, requestTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
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
		Dial: timeoutDialer(time.Duration(connectTimeout)*time.Second, time.Duration(requestTimeout)*time.Second),
	}
	return xmlrpc.NewClient(url, &transport)
}
