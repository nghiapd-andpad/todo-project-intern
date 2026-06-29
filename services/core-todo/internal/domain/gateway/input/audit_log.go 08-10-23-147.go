package input

type CreateAuditLog struct {
	ActorID    int64
	EventName  string
	EntityType string
	EntityID   int64
	Payload    []byte
}
