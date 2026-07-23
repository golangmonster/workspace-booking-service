package booking

import (
	"context"

	pb "github.com/golangmonster/workspace-booking-service/pkg/api/booking/v1"
)

func (i *Implementation) ListBookings(ctx context.Context, req *pb.ListBookingsRequest) (*pb.ListBookingsResponse, error) {
	resp, err := i.bookingService.ListBookings(ctx, toListBookingsRequest(req))
	if err != nil {
		return nil, err
	}

	return toListBookingsResponse(resp), nil
}
