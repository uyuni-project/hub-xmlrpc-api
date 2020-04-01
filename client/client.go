package client

import (
	"net"
	"net/http"
	"time"

	"github.com/kolo/xmlrpc"
)

type Client struct {
	connectTimeout, readWriteTimeout int
}

func NewClient(connectTimeout, readWriteTimeout int) *Client {
	return &Client{connectTimeout: connectTimeout, readWriteTimeout: readWriteTimeout}
}

func (c *Client) ExecuteCall(url string, call string, args []interface{}) (response interface{}, err error) {
	client, err := getClientWithTimeout(url, c.connectTimeout, c.readWriteTimeout)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	err = client.Call(call, args, &response)
	return response, err
}

func timeoutDialer(connectTimeout, readWriteTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, connectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(readWriteTimeout))
		return conn, nil
	}
}

func getClientWithTimeout(url string, connectTimeout, readWriteTimeout int) (*xmlrpc.Client, error) {
	transport := http.Transport{
		Dial: timeoutDialer(time.Duration(connectTimeout)*time.Second, time.Duration(readWriteTimeout)*time.Second),
	}
	return xmlrpc.NewClient(url, &transport)
}
