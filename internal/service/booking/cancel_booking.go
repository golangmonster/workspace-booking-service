package booking

import (
	"context"

	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
)

func (s *service) CancelBooking(ctx context.Context, bookingID int64) error {
	err := s.repo.UpdateBookingStatus(ctx, UpdateBookingStatus{
		BookingID: bookingID,
		Status:    booking.StatusCancelled,
	})
	if err != nil {
		return err
	}

	return nil
}
