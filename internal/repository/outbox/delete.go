package outbox

import (
	"context"

	"github.com/Masterminds/squirrel"
)

func (r *repository) DeleteOutboxItems(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	qb := squirrel.Delete("outbox").
		Where(
			squirrel.Eq{"id": ids},
		)

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Querier(ctx).Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
