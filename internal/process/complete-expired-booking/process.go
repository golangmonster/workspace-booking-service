package complete_expired_booking

import (
	"context"
	"fmt"
	"time"
)

type process struct {
	bookingRepo bookingRepository
}

func NewProcess(bookingRepo bookingRepository) *process {
	return &process{
		bookingRepo: bookingRepo,
	}
}

func (p *process) Run(ctx context.Context) error {
	err := p.bookingRepo.CompleteExpiredBookings(ctx, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("complete expired bookings: %w", err)
	}

	return nil
}
