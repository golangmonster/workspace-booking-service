package outbox

import (
	"context"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

const (
	limit = 100
)

type process struct {
	producer   producer
	outboxRepo outboxRepository
	topic      string
}

func NewProcess(producer producer, outboxRepo outboxRepository, topic string) *process {
	return &process{
		producer:   producer,
		outboxRepo: outboxRepo,
		topic:      topic,
	}
}

func (p *process) Run(ctx context.Context) error {
	items, err := p.outboxRepo.ListOutboxItems(ctx, p.topic, limit)
	if err != nil {
		return fmt.Errorf("list outbox items: %w", err)
	}

	// Сортировка по возрастанию ключа
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	deleteIDs := make([]int64, 0)

	for _, item := range items {
		err = p.producer.Produce(p.topic, item.Value, item.AggregateID)
		if err != nil {
			log.Error(err, "failed to send kafka message for topic ", p.topic)

			break
		}

		deleteIDs = append(deleteIDs, item.ID)
	}

	err = p.outboxRepo.DeleteOutboxItems(ctx, deleteIDs)
	if err != nil {
		return fmt.Errorf("delete outbox items: %w", err)
	}

	return nil
}
