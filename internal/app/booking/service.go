package booking

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/golangmonster/workspace-booking-service/pkg/api/booking/v1"
)

type Implementation struct {
	pb.UnimplementedBookingServiceServer

	bookingService bookingService
}

func New(service bookingService) *Implementation {
	return &Implementation{
		bookingService: service,
	}
}

func (i *Implementation) RegisterServer(server *grpc.Server) {
	pb.RegisterBookingServiceServer(server, i)
}

func (i *Implementation) RegisterHandlerFromEndpoint(
	ctx context.Context,
	mux *runtime.ServeMux,
	addrGRPC string,
	opts []grpc.DialOption,
) error {
	return pb.RegisterBookingServiceHandlerFromEndpoint(ctx, mux, addrGRPC, opts)
}
