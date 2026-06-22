package event

import "time"

type NotificationEmailRequested struct {
	NotificationID int64     `json:"notification_id"`
	ReceiverID     int64     `json:"receiver_id"`
	Email          string    `json:"email"`
	Subject        string    `json:"subject"`
	Content        string    `json:"content"`
	OccurredOn     time.Time `json:"occurred_at"`
}

func (e NotificationEmailRequested) EventName() string {
	return "notification.email.requested"
}

func (e NotificationEmailRequested) OccurredAt() time.Time {
	return e.OccurredOn
}
