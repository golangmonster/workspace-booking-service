package outbox

import (
	"context"

	"github.com/golangmonster/workspace-booking-service/internal/model/outbox"
)

type producer interface {
	Produce(topic string, msg string, key string) error
}

type outboxRepository interface {
	DeleteOutboxItems(ctx context.Context, ids []int64) error
	ListOutboxItems(ctx context.Context, topic string, limit uint64) ([]outbox.Outbox, error)
}
