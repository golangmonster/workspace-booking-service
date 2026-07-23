package booking

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
)

func (r *repository) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (int64, error) {
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
		Suffix("RETURNING id")

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	var id int64

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &id, sql, args...)
	if err != nil {
		return 0, err
	}

	return id, nil
}
