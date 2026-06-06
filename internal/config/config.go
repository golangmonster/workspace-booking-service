package config

type Config struct {
	HTTPAddr string `envconfig:"HTTP_ADDR" default:":8000"`
	GRPCAddr string `envconfig:"GRPC_ADDR" default:":8001"`

	PostgresDSN string `envconfig:"POSTGRES_DSN" required:"true"`
}
