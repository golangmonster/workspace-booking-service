package booking

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
)

func (r *repository) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (*booking.Booking, error) {
	t := time.Now().UTC()

	qb := squirrel.Insert("booking").
		Columns(
			"user_id",
			"workspace_id",
			"start_at",
			"end_at",
			"status",
			"created_at",
			"updated_at",
		).
		Values(
			req.UserID,
			req.WorkspaceID,
			req.StartAt,
			req.EndAt,
			req.Status,
			t,
			t,
		).
		Suffix("RETURNING id, user_id, workspace_id, start_at, end_at, status, created_at, updated_at")

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var item bookingItem

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		return nil, err
	}

	return toBooking(item), nil
}
