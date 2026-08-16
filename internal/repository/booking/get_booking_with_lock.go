package booking

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
)

func (r *repository) GetBookingWithLock(ctx context.Context, bookingID int64) (*booking.Booking, error) {
	qb := squirrel.Select(
		"id",
		"workspace_id",
		"user_id",
		"start_at",
		"end_at",
		"status",
		"created_at",
		"updated_at",
	).From("booking").
		Where(squirrel.Eq{"id": bookingID}).
		PlaceholderFormat(squirrel.Dollar)

	var item bookingItem

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		return nil, err
	}

	return &booking.Booking{
		ID:          item.ID,
		WorkspaceID: item.WorkspaceID,
		UserID:      item.UserID,
		StartAt:     item.StartAt,
		EndAt:       item.EndAt,
		Status:      booking.Status(item.Status),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}
