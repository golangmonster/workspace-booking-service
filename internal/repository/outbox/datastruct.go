package outbox

type outboxItem struct {
	ID           int64  `db:"id"`
	Topic        string `db:"topic"`
	AggregateID  string `db:"aggregate_id"`
	MessageValue string `db:"message_value"`
}
