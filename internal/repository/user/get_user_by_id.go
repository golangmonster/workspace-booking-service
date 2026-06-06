package user

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/user"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	qb := squirrel.Select(
		"id",
		"login",
		"full_name",
		"phone",
		"created_at",
		"updated_at",
	).From("users").
		Where(squirrel.Eq{
			"id":         id,
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

func (r *repository) GetUsersByIDs(ctx context.Context, ids ...int64) ([]*user.User, error) {
	qb := squirrel.Select(
		"id",
		"login",
		"full_name",
		"phone",
		"created_at",
		"updated_at",
	).From("users").
		Where(squirrel.Eq{
			"id":         ids,
			"is_deleted": false,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var items []userItem

	err = pgxscan.Select(ctx, r.pool.Querier(ctx), &items, sql, args...)
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, 0, len(ids))
	for _, item := range items {
		users = append(users, toUser(item))
	}

	return users, nil
}
