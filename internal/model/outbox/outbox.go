package outbox

type Message struct {
	Topic       string
	AggregateID string
	Value       string
}

type Outbox struct {
	Message

	ID int64
}
