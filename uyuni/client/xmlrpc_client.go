package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/http"
	neturl "net/url"
	"sync"
	"time"

	xmlrpc "github.com/uyuni-project/xmlrpc-public-methods"
)

type Client struct {
	connectTimeout, requestTimeout int
	rootCAs                        *x509.CertPool
	mu                             sync.RWMutex
	trace                          bool
}

func NewClient(connectTimeout, requestTimeout int, trace bool) *Client {
	return &Client{
		connectTimeout: connectTimeout,
		requestTimeout: requestTimeout,
		rootCAs:        x509.NewCertPool(),
		trace:          trace,
	}
}

func (c *Client) AddRootCA(pem []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.trace {
		log.Printf("Adding root CA:\n %s\n", string(pem))
	}
	if !c.rootCAs.AppendCertsFromPEM(pem) {
		return fmt.Errorf("failed to parse root CA certificate")
	}
	return nil
}

func (c *Client) ExecuteCall(endpoint string, call string, args []interface{}) (response interface{}, err error) {
	if c.trace {
		log.Printf("Calling %s on %s\n", call, endpoint)
	}
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
	tlsConfig := tls.Config{
		RootCAs: c.rootCAs,
	}
	if c.trace {
		tlsConfig.InsecureSkipVerify = true // MUST be true, otherwise VerifyPeerCertificate is skipped on failure
		// Hook to trace SSL:
		tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			fmt.Println("\n=== TLS HANDSHAKE TRACE ===")
			fmt.Printf("Server sent %d certificates in its chain:\n", len(rawCerts))

			certs := make([]*x509.Certificate, len(rawCerts))
			for i, asn1Data := range rawCerts {
				cert, err := x509.ParseCertificate(asn1Data)
				if err != nil {
					fmt.Printf("  [%d] Failed to parse raw cert: %v\n", i, err)
					return err
				}

				pemBlock := &pem.Block{
					Type:  "CERTIFICATE",
					Bytes: asn1Data,
				}
				pemBytes := pem.EncodeToMemory(pemBlock)

				certs[i] = cert
				fmt.Printf("  [%d] Subject: % s\n      Issuer:  %s\n%s\n", i, cert.Subject, cert.Issuer, string(pemBytes))
			}

			opts := x509.VerifyOptions{
				Roots:         c.rootCAs,
				CurrentTime:   time.Now(),
				Intermediates: x509.NewCertPool(),
			}

			for _, cert := range certs[1:] {
				opts.Intermediates.AddCert(cert)
			}

			_, err := certs[0].Verify(opts)
			if err != nil {
				fmt.Printf("\n[!] Cryptographic Validation Error: %v\n", err)
			} else {
				fmt.Println("[-] Cryptographic verification inside hook passed!")
			}
			if u, uErr := neturl.Parse(url); uErr == nil {
				if hErr := certs[0].VerifyHostname(u.Hostname()); hErr != nil {
					fmt.Printf("[!] Hostname Validation Error: %v\n", hErr)
				} else {
					fmt.Printf("[-] Hostname verification passed for %s\n", u.Hostname())
				}
			}
			if len(certs[0].IssuingCertificateURL) > 0 {
				fmt.Println("[WARNING] This certificate utilizes Authority Information Access (AIA).")
				fmt.Printf("          Missing intermediates might need to be downloaded from: %v\n", certs[0].IssuingCertificateURL)
			}
			fmt.Println("===========================")

			return err
		}
	}
	transport := http.Transport{
		DialContext:     timeoutDialer(time.Duration(connectTimeout)*time.Second, time.Duration(requestTimeout)*time.Second),
		TLSClientConfig: &tlsConfig,
	}
	return xmlrpc.NewClient(url, &transport)
}
