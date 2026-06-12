package input

type CreateOutboxEvent struct {
	EventName  string
	RoutingKey string
	Payload    []byte
}

type ListPendingOutboxEvents struct {
	Limit int
}
