package outbox

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golangmonster/workspace-booking-service/internal/model/outbox"
	"github.com/samber/lo"
)

func (r *repository) ListOutboxItems(ctx context.Context, topic string, limit uint64) ([]outbox.Outbox, error) {
	qb := squirrel.Select(
		"id",
		"topic",
		"aggregate_id",
		"message_value",
	).From("outbox").
		Where(squirrel.Eq{"topic": topic}).
		OrderBy("id").
		Limit(limit)

	sql, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var items []outboxItem
	
	if err := pgxscan.Select(ctx, r.pool.Querier(ctx), &items, sql, args...); err != nil {
		return nil, err
	}

	return lo.Map(items, func(i outboxItem, _ int) outbox.Outbox {
		return outbox.Outbox{
			ID: i.ID,
			Message: outbox.Message{
				AggregateID: i.AggregateID,
				Topic:       i.Topic,
				Value:       i.MessageValue,
			},
		}
	}), nil
}
