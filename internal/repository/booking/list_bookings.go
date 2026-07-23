package booking

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/booking"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/booking"
	"github.com/samber/lo"
)

func (r *repository) ListBookings(ctx context.Context, req dto.ListBookingsRequest) (dto.ListBookingsResponse, error) {
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
		PlaceholderFormat(squirrel.Dollar)

	qb = bookingFilter(qb, req.Filter)

	if req.Page != nil {
		qb = qb.Limit(req.Page.Limit).Offset(req.Page.Offset)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return dto.ListBookingsResponse{}, err
	}

	var items []bookingItem

	err = pgxscan.Select(ctx, r.pool.Querier(ctx), &items, sql, args...)
	if err != nil {
		return dto.ListBookingsResponse{}, err
	}

	totalCount, err := r.getTotalBookingsCount(ctx, req.Filter)
	if err != nil {
		return dto.ListBookingsResponse{}, err
	}

	return dto.ListBookingsResponse{
		Bookings: lo.Map(items, func(item bookingItem, _ int) *booking.Booking {
			return toBooking(item)
		}),
		TotalCount: totalCount,
	}, nil
}

func (r *repository) getTotalBookingsCount(ctx context.Context, filter *dto.Filter) (uint32, error) {
	qb := squirrel.Select("COUNT(id)").
		From("booking").
		PlaceholderFormat(squirrel.Dollar)

	qb = bookingFilter(qb, filter)

	sql, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var count uint32

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &count, sql, args...)

	return count, err
}

func bookingFilter(qb squirrel.SelectBuilder, f *dto.Filter) squirrel.SelectBuilder {
	if f == nil {
		return qb
	}

	if f.UserID != nil {
		qb = qb.Where(squirrel.Eq{"user_id": *f.UserID})
	}

	if f.WorkspaceID != nil {
		qb = qb.Where(squirrel.Eq{"workspace_id": *f.WorkspaceID})
	}

	if len(f.Statuses) > 0 {
		qb = qb.Where(squirrel.Eq{"status": f.Statuses})
	}

	// Overlap check: a booking [start_at, end_at) is within the requested
	// range of days if it ends after the start of DateFrom and starts
	// before the day after DateTo.
	if f.DateFrom != nil {
		qb = qb.Where(squirrel.Gt{"end_at": startOfDay(*f.DateFrom)})
	}

	if f.DateTo != nil {
		qb = qb.Where(squirrel.Lt{"start_at": startOfDay(*f.DateTo).AddDate(0, 0, 1)})
	}

	return qb
}

func startOfDay(t time.Time) time.Time {
	t = t.UTC()

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
