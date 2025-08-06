package config

import (
	"crypto/tls"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/middlewares"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func (c *config) InitializeGRPC() {

	var err error
	c.listener, err = net.Listen(c.env.GRPC.Network, c.env.GRPC.Port)
	if err != nil {
		slog.Error("failed to listen", " error", err)
		panic(err)
	}

	slog.Info("Server listening on", "Addr:", c.listener.Addr())

	// Load server's certificate and private key
	homeDir, _ := os.UserHomeDir()
	certPath := filepath.Join(homeDir, "grpc-ssl", "fullchain.pem")
	keyPath := filepath.Join(homeDir, "grpc-ssl", "privkey.pem")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		slog.Error("failed to load key pair", "error", err)
		panic(err)
	}

	// Create a new gRPC server with the telemetry interceptor and TLS credentials
	c.server = grpc.NewServer(
		grpc.Creds(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
		})),
		grpc.ChainUnaryInterceptor(
			middlewares.TelemetryInterceptor(c.context),
			middlewares.AuthInterceptor(c.context, c.activity),
			middlewares.AccessControlInterceptor(c.context, &c.cache),
		),
	)

	reflection.Register(c.server)

}

func (c *config) GetGRPCServer() *grpc.Server {
	return c.server
}

func (c *config) GetListener() net.Listener {
	return c.listener
}

func (c *config) GetInfos() (serviceQty int, methodQty int) {
	serviceQty = len(c.server.GetServiceInfo())
	for _, s := range c.server.GetServiceInfo() {
		methodQty += len(s.Methods)
	}
	return
}
