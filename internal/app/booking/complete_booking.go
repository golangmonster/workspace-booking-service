package booking

import (
	"context"
	"errors"

	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/booking/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CompleteBooking(ctx context.Context, req *pb.CompleteBookingRequest) (*pb.CompleteBookingResponse, error) {
	err := i.bookingService.UpdateBookingStatus(ctx, dto.UpdateBookingStatus{
		BookingID: req.Id,
		Status:    booking.StatusCompleted,
	})
	if err != nil {
		switch {
		case errors.Is(err, booking.ErrBookingNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, booking.ErrInvalidStatusTransition):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, err
	}

	return &pb.CompleteBookingResponse{}, nil
}
