package config

import "time"

type Config struct {
	HTTPAddr string `envconfig:"HTTP_ADDR" default:":8000"`
	GRPCAddr string `envconfig:"GRPC_ADDR" default:":8001"`

	PostgresDSN string `envconfig:"POSTGRES_DSN" required:"true"`

	CompleteExpiredBookingDuration time.Duration `envconfig:"COMPLETE_EXPIRED_BOOKING_DURATION" default:"10s"`
	CompleteExpiredBookingEnabled  bool          `envconfig:"COMPLETE_EXPIRED_BOOKING_ENABLED" default:"false"`
}
