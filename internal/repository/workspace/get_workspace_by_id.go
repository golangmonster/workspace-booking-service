package workspace

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
)

func (r *repository) GetWorkspaceByID(ctx context.Context, id int64) (*model.Workspace, error) {
	qb := squirrel.Select(
		"id",
		"name",
		"lat",
		"lon",
		"full_address",
		"type",
		"status",
		"capacity",
		"created_at",
		"updated_at",
		"is_deleted",
	).From("workspace").
		Where(squirrel.Eq{
			"id":         id,
			"is_deleted": false,
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var item workspaceItem

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrWorkspaceNotFound
		}

		return nil, err
	}

	return toWorkspace(item), nil
}
