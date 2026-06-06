package user

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golangmonster/workspace-booking-service/internal/model/user"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/user"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) error {
	ut := time.Now().UTC()

	qb := squirrel.Update("users").
		Where(squirrel.Eq{"id": req.ID, "is_deleted": false}).
		Set("login", req.Login).
		Set("phone", req.Phone).
		Set("full_name", req.FullName).
		Set("updated_at", ut).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	tag, err := r.pool.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				switch pgErr.ConstraintName {
				case loginUniqueConstraintName:
					return user.ErrLoginExists
				case phoneUniqueConstraintName:
					return user.ErrPhoneExists
				}
			}
		}

		return err
	}

	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
