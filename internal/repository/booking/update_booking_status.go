package booking

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
)

func (r *repository) UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatus) error {
	t := time.Now().UTC()

	qb := squirrel.Update("booking").
		Set("status", req.Status).
		Set("updated_at", t).
		Where(squirrel.Eq{"id": req.BookingID})

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	tag, err := r.pool.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return booking.ErrBookingNotFound
	}

	return nil
}
