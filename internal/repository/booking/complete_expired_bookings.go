package booking

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
)

func (r *repository) CompleteExpiredBookings(ctx context.Context, endAtLt time.Time) error {
	qb := squirrel.Update("booking").
		Set("status", booking.StatusCompleted).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.And{
			squirrel.Lt{"end_at": endAtLt},
			squirrel.Eq{"status": booking.StatusActive},
		})

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
