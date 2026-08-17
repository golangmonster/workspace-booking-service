package complete_expired_booking

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

type process struct {
	bookingRepo bookingRepository
}

func NewProcess(bookingRepo bookingRepository) *process {
	return &process{
		bookingRepo: bookingRepo,
	}
}

func (p *process) Run(ctx context.Context) {
	go func() {
		err := p.bookingRepo.CompleteExpiredBookings(ctx, time.Now().UTC())
		if err != nil {
			log.Error("complete expired bookings: ", err)
		}
	}()
}
