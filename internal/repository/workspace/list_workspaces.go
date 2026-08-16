package workspace

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/workspace"
	"github.com/samber/lo"
)

func (r *repository) ListWorkspaces(ctx context.Context, req *dto.ListWorkspacesRequest) (*dto.ListWorkspacesResponse, error) {
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
		Where(squirrel.Eq{"is_deleted": false}).
		PlaceholderFormat(squirrel.Dollar)

	qb = workspaceFilter(qb, req.Filter)

	if req.Page != nil {
		qb = qb.Offset(req.Page.Offset).Limit(req.Page.Limit)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var items []workspaceItem

	err = pgxscan.Select(ctx, r.pool.Querier(ctx), &items, sql, args...)
	if err != nil {
		return nil, err
	}

	totalCount, err := r.getTotalWorkspacesCount(ctx, req.Filter)
	if err != nil {
		return nil, err
	}

	return &dto.ListWorkspacesResponse{
		Workspaces: lo.Map(items, func(i workspaceItem, _ int) *model.Workspace {
			return toWorkspace(i)
		}),
		TotalCount: totalCount,
	}, nil
}

func workspaceFilter(qb squirrel.SelectBuilder, f *dto.WorkspaceFilter) squirrel.SelectBuilder {
	if f == nil {
		return qb
	}

	if len(f.Types) > 0 {
		qb = qb.Where(squirrel.Eq{"type": f.Types})
	}

	if f.CapacityGte != nil {
		qb = qb.Where(squirrel.GtOrEq{"capacity": *f.CapacityGte})
	}

	if f.FullAddress != nil {
		qb = qb.Where(squirrel.ILike{"full_address": "%" + *f.FullAddress + "%"})
	}

	return qb
}

func (r *repository) getTotalWorkspacesCount(ctx context.Context, filter *dto.WorkspaceFilter) (uint32, error) {
	qb := squirrel.Select("COUNT(id)").From("workspace")

	qb = workspaceFilter(qb, filter)

	sql, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var count uint32

	err = pgxscan.Get(ctx, r.pool.Querier(ctx), &count, sql, args...)

	return count, err
}
