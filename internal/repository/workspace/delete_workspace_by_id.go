package workspace

import (
	"context"

	"github.com/Masterminds/squirrel"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
)

func (r *repository) DeleteWorkspaceByID(ctx context.Context, id int64) error {
	qb := squirrel.Update("workspace").
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
		return model.ErrWorkspaceNotFound
	}

	return nil
}
