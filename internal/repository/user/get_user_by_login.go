package user

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/user"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	qb := squirrel.Select(
		"id",
		"login",
		"full_name",
		"phone",
		"created_at",
		"updated_at",
	).From("users").
		Where(squirrel.Eq{
			"login":      login,
			"is_deleted": false,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var item userItem

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}

		return nil, err
	}

	return toUser(item), nil
}
