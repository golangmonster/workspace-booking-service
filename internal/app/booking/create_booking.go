package booking

import (
	"context"
	"errors"

	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	"github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/booking/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	id, err := i.bookingService.CreateBooking(ctx, dto.CreateBookingRequest{
		UserID:      req.UserId,
		WorkspaceID: req.WorkspaceId,
		Status:      booking.StatusActive,
		StartAt:     req.StartAt.AsTime(),
		EndAt:       req.EndAt.AsTime(),
	})
	if err != nil {
		switch {
		case errors.Is(err, booking.ErrInvalidTimeRange):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, workspace.ErrWorkspaceNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, workspace.ErrWorkspaceNotAvailable):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, booking.ErrBookingOverlap):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, err
	}

	return &pb.CreateBookingResponse{
		BookingId: id,
	}, nil
}
