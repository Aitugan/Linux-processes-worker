package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Port          int
	Host          string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	TrustStore    string
	ServerKey     string
	ServerCert    string
	ServerCaCerts string
}

type Server struct {
	config *ServerConfig
	server *http.Server
}

func NewServer(c *ServerConfig) (*Server, error) {
	s := &Server{config: c}
	server := &http.Server{
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
		Addr:         fmt.Sprintf(":%d", s.config.Port),
	}

	caCertPool := x509.NewCertPool()

	trustStore, err := ioutil.ReadFile(s.config.TrustStore)
	if err != nil {
		return nil, err
	}
	caCertPool.AppendCertsFromPEM(trustStore)

	serverCert, err := ioutil.ReadFile(s.config.ServerCert)
	if err != nil {
		return nil, err
	}
	var certs []byte
	certs = append(certs, []byte(serverCert)...)

	serverCACerts, err := ioutil.ReadFile(s.config.ServerCaCerts)
	if err != nil {
		return nil, err
	}
	certPEMBlock := []byte(serverCACerts)
	var certDERBlock *pem.Block
	isRootCert := true
	for {
		// Extract all intermedidate certificates
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if isRootCert {
			isRootCert = false
			continue
		}

		if certDERBlock.Type == "CERTIFICATE" {
			cert := &bytes.Buffer{}
			pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: certDERBlock.Bytes})
			certs = append(certs, cert.Bytes()...)
		}
	}

	serverKey, err := ioutil.ReadFile(s.config.ServerKey)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(certs, serverKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server keypair. Error: %s", err)
	}

	tlsConfig := &tls.Config{
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}
	tlsConfig.BuildNameToCertificate()
	server.TLSConfig = tlsConfig

	router, err := New()
	if err != nil {
		return nil, err
	}
	server.Handler = router

	return &Server{config: c, server: server}, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Start() error {
	log.Info("started server on Port: ", s.server.Addr)
	err := s.server.ListenAndServeTLS("", "")
	if err != nil {
		return err
	}
	return nil
}
