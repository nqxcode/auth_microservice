package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
	grpcCertEnvName = "GRPC_CERT"
	grpcKeyEnvName  = "GRPC_KEY"
)

// GRPCConfig grpc server config
type GRPCConfig interface {
	Address() string
	Cert() []byte
	Key() []byte
}

type grpcConfig struct {
	host string
	port string
	cert []byte
	key  []byte
}

// NewGRPCConfig create new grpc server config
func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	cert := os.Getenv(grpcCertEnvName)
	if len(cert) == 0 {
		return nil, errors.New("grpc cert not found")
	}

	key := os.Getenv(grpcKeyEnvName)
	if len(key) == 0 {
		return nil, errors.New("grpc key not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
		cert: []byte(cert),
		key:  []byte(key),
	}, nil
}

// Address get grpc server address
func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

// Cert get service pem cert
func (cfg *grpcConfig) Cert() []byte {
	return cfg.cert
}

// Key get service key
func (cfg *grpcConfig) Key() []byte {
	return cfg.key
}
