package booking

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
	"github.com/jackc/pgx/v5"
)

func (r *repository) UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatus) (*booking.Booking, error) {
	t := time.Now().UTC()

	qb := squirrel.Update("booking").
		Set("status", req.Status).
		Set("updated_at", t).
		Where(squirrel.Eq{"id": req.BookingID}).
		Suffix("RETURNING id, user_id, workspace_id, start_at, end_at, status, created_at, updated_at")

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var item bookingItem

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, booking.ErrBookingNotFound
		}

		return nil, err
	}

	return toBooking(item), nil
}
