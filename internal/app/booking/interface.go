package booking

import (
	"context"

	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
)

type bookingService interface {
	CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (int64, error)
	ListBookings(ctx context.Context, req dto.ListBookingsRequest) (dto.ListBookingsResponse, error)
	UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatus) error
}
