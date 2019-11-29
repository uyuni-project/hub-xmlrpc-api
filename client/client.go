package client

import (
	"net"
	"net/http"
	"time"

	"github.com/kolo/xmlrpc"
)

var defaultConnectTimeout, defaultReadWriteTimeout = 1, 1

type timeoutConfig struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func timeoutDialer(config *timeoutConfig) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(config.ReadWriteTimeout))
		return conn, nil
	}
}

func GetClientWithTimeout(url string, connectTimeout, readWriteTimeout int) (*xmlrpc.Client, error) {
	timeoutConfig := &timeoutConfig{
		ConnectTimeout:   time.Duration(connectTimeout) * time.Second,
		ReadWriteTimeout: time.Duration(readWriteTimeout) * time.Second,
	}

	transport := http.Transport{
		Dial: timeoutDialer(timeoutConfig),
	}
	return xmlrpc.NewClient(url, &transport)
}

func GetDefaultTimeoutClient(url string) (*xmlrpc.Client, error) {
	return GetClientWithTimeout(url, defaultConnectTimeout, defaultReadWriteTimeout)
}
