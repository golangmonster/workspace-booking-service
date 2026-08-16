package booking

import (
	"context"
	"fmt"

	model "github.com/golangmonster/workspace-booking-service/internal/model/booking"
)

func (s *service) UpdateBookingStatus(ctx context.Context, req UpdateBookingStatus) error {
	return s.repo.InTx(ctx, func(ctx context.Context) error {
		booking, err := s.repo.GetBookingWithLock(ctx, req.BookingID)
		if err != nil {
			return fmt.Errorf("get booking with lock: %w", err)
		}

		if !model.StatusTransitionAllowed(booking.Status, req.Status) {
			return fmt.Errorf("%w from %s to %s", model.ErrInvalidStatusTransition, booking.Status, req.Status)
		}

		err = s.repo.UpdateBookingStatus(ctx, req)
		if err != nil {
			return fmt.Errorf("update booking status: %w", err)
		}

		return nil
	})
}
