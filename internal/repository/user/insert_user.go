package user

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/user"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/user"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) InsertUser(ctx context.Context, req *dto.CreateUserRequest) (*user.User, error) {
	ct := time.Now().UTC()

	qb := squirrel.Insert("users").Columns(
		"login",
		"full_name",
		"phone",
		"created_at",
		"updated_at",
		"is_deleted",
	).Values(
		req.Login,
		req.FullName,
		req.Phone,
		ct,
		ct,
		false,
	).Suffix(`RETURNING id, login, full_name, 
	phone, created_at, updated_at, is_deleted`).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var item userItem

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				switch pgErr.ConstraintName {
				case loginUniqueConstraintName:
					return nil, user.ErrLoginExists
				case phoneUniqueConstraintName:
					return nil, user.ErrPhoneExists
				}
			}
		}

		return nil, err
	}

	return toUser(item), nil
}
