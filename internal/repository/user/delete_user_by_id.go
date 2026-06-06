package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/golangmonster/workspace-booking-service/internal/model/user"
)

func (r *repository) DeleteUserByID(ctx context.Context, id int64) error {
	qb := squirrel.Update("users").
		Set("is_deleted", true).
		Where(squirrel.Eq{"id": id, "is_deleted": false}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	tag, err := r.pool.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
