package workspace

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/workspace"
)

func (r *repository) UpdateWorkspace(ctx context.Context, req *dto.UpdateWorkspaceRequest) error {
	ut := time.Now().UTC()

	qb := squirrel.Update("workspace").
		Where(squirrel.Eq{"id": req.ID, "is_deleted": false}).
		Set("name", req.Name).
		Set("lat", req.Lat).
		Set("lon", req.Lon).
		Set("full_address", req.FullAddress).
		Set("capacity", req.Capacity).
		Set("status", req.Status).
		Set("updated_at", ut).
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
