# Event Driven Flow

## Logic

CreateTodo API
--> TodoCreator xử lý business logic
--> INSERT todos
--> INSERT outbox_events với status = PENDING
--> Commit DB transaction

Outbox Publisher Worker
--> định kỳ query outbox_events có status PENDING/FAILED
--> publish event vào RabbitMQ exchange `todo.events`
--> RabbitMQ route theo routing key, ví dụ `todo.created`
--> message được đưa vào queue `todo.audit.queue`
--> Audit Consumer consume message
--> INSERT audit_logs
--> ACK message
--> RabbitMQ xoá message khỏi queue

## Flow
CreateTodo API (core-todo)
      |
      v
TodoCreator
      |
      |-- INSERT todos
      |-- INSERT outbox_events status = PENDING
      v
DB Transaction Commit

------------------------

Outbox Publisher Worker (core-user)
      |
      |-- SELECT PENDING / FAILED events (limit)
      |-- Publish to RabbitMQ
      v
RabbitMQ Exchange: todo.events
      |
      | routing_key = todo.assigned
      v
RabbitMQ Queue: todo.notification.queue (binding)
      |
      v
Notification Consumer
      |
      |-- INSERT notifications tables, gửi email cho user
      |-- ACK message
      v
RabbitMQ deletes message

## Current Progress

Done:
- Added Domain Event definitions.
- Added `outbox_events` and `audit_logs` schemas/migrations.
- Added outbox/audit persistence layer.
- Added RabbitMQ config, publisher, and topology setup.
- Integrated `TodoCreated` event into `CreateTodo`.

Next:
- Implement Outbox Publisher Worker.
- Implement Audit Consumer.
- Apply events to remaining Todo/TodoList CRUD operations.


## Improve notification consumer
NotificationConsumer nhận todo.assigned
    │
    ├── INSERT notifications (DB)
    │
    └── INSERT outbox_events (notification.created)
              │
              ▼
    OutboxPublisherJob (core-user worker)
              │
              ▼
    user.events exchange
              │
              ▼
    email.queue
              │
              ▼
    EmailConsumer → SMTP