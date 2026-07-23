package booking

import (
	"context"
	"errors"

	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/booking/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	err := i.bookingService.CancelBooking(ctx, req.Id)
	if err != nil {
		if errors.Is(err, booking.ErrBookingNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &pb.CancelBookingResponse{}, nil
}
