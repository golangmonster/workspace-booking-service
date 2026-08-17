package complete_expired_booking

import (
	"context"
	"time"
)

type bookingRepository interface {
	CompleteExpiredBookings(ctx context.Context, endAtLt time.Time) error
}
