package outbox

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golangmonster/workspace-booking-service/internal/model/outbox"
)

func (r *repository) CreateOutboxItem(ctx context.Context, item outbox.Message) error {
	ct := time.Now().UTC()

	qb := squirrel.Insert("outbox").
		Columns(
			"topic",
			"aggregate_id",
			"message_value",
			"created_at",
		).
		Values(
			item.Topic,
			item.AggregateID,
			item.Value,
			ct,
		)

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	if _, err := r.pool.Querier(ctx).Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}
