package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aitugan/CodingChallenge/internal/api"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var startServerCmd = &cobra.Command{
	Use:   "startServer",
	Short: "A start command",
	Long:  "A start command",
	Run: func(cmd *cobra.Command, args []string) {

		sc := &api.ServerConfig{
			Port:          8080,
			Host:          "localhost",
			ReadTimeout:   5 * time.Second,
			WriteTimeout:  5 * time.Second,
			IdleTimeout:   10 * time.Second,
			TrustStore:    "./internal/api/certs/Client-CA-cert.pem",
			ServerKey:     "./internal/api/certs/server-key.pem",
			ServerCert:    "./internal/api/certs/server-cert.pem",
			ServerCaCerts: "./internal/api/certs/CA-cert.pem",
		}
		server, err := api.NewServer(sc)
		if err != nil {
			log.Fatal("failed to create server, Error:", err)
		}
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()

		sig := make(chan os.Signal, 0)
		signal.Notify(sig,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		go func() {
			err := server.Start()
			if err != nil {
				log.Error("server error ", err)
			}
		}()
		log.Info("server started and ready accept requests and signals")

		select {
		case <-sig:
			log.Info("received close signal")
			server.Shutdown(ctxShutDown)
		}

	},
}
