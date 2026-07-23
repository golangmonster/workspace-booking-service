package booking

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	"github.com/jackc/pgx/v5"
)

type workspaceItem struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Lat         float64   `db:"lat"`
	Lon         float64   `db:"lon"`
	FullAddress string    `db:"full_address"`
	Type        string    `db:"type"`
	Status      string    `db:"status"`
	Capacity    uint32    `db:"capacity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	IsDeleted   bool      `db:"is_deleted"`
}

// GetWorkspaceForUpdate locks the workspace row (SELECT ... FOR UPDATE) to serialize
// concurrent booking attempts for the same workspace, and returns it so the caller
// can validate its status.
func (r *repository) GetWorkspaceForUpdate(ctx context.Context, workspaceID int64) (*workspace.Workspace, error) {
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
		Where(squirrel.Eq{"id": workspaceID}).
		Suffix("FOR UPDATE").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var item workspaceItem

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &item, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, workspace.ErrWorkspaceNotFound
		}

		return nil, err
	}

	return &workspace.Workspace{
		ID:          item.ID,
		Name:        item.Name,
		Lat:         item.Lat,
		Lon:         item.Lon,
		FullAddress: item.FullAddress,
		Type:        workspace.Type(item.Type),
		Status:      workspace.Status(item.Status),
		Capacity:    item.Capacity,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		IsDeleted:   item.IsDeleted,
	}, nil
}
