package grpc

import (
	govalidator "buf.build/go/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func WithValidation() grpc.UnaryServerInterceptor {
	v, err := govalidator.New()
	if err != nil {
		log.Fatal("failed to initialize validator ", err)
	}

	return protovalidate.UnaryServerInterceptor(v)
}
