package user

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/golangmonster/workspace-booking-service/pkg/api/user/v1"
)

type Implementation struct {
	pb.UnimplementedUserServiceServer

	userService userService
}

func New(service userService) *Implementation {
	return &Implementation{
		userService: service,
	}
}

func (i *Implementation) RegisterServer(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, i)
}

func (i *Implementation) RegisterHandlerFromEndpoint(
	ctx context.Context,
	mux *runtime.ServeMux,
	addrGRPC string,
	opts []grpc.DialOption,
) error {
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, addrGRPC, opts)

	return err
}
