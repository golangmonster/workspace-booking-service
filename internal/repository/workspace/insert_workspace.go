package workspace

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"

	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/workspace"
)

func (r *repository) InsertWorkspace(ctx context.Context, req *dto.CreateWorkspaceRequest) (*model.Workspace, error) {
	t := time.Now().UTC()

	qb := squirrel.Insert("workspace").Columns(
		"name",
		"lat",
		"lon",
		"full_address",
		"type",
		"status",
		"capacity",
		"created_at",
		"updated_at",
	).Values(
		req.Name,
		req.Lat,
		req.Lon,
		req.FullAddress,
		req.Type,
		req.Status,
		req.Capacity,
		t,
		t,
	).Suffix(`RETURNING id, name, lat, lon, full_address,
	type, status, capacity, created_at, updated_at`).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var item workspaceItem

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				switch pgErr.ConstraintName {
				case workspaceNameTypeUniqueConstraint:

					return nil, model.ErrWorkspaceAlreadyExists
				}
			}
		}

		return nil, err
	}

	return toWorkspace(item), nil
}
