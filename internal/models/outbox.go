package models

import "time"

// OutboxEvent represents an event stored in Postgres outbox_events table
type OutboxEvent struct {
    ID             int64                  `db:"id" json:"id"`
    AggregateType  string                 `db:"aggregate_type" json:"aggregate_type"`
    AggregateID    string                 `db:"aggregate_id" json:"aggregate_id"`
    EventType      string                 `db:"event_type" json:"event_type"`
    IdempotencyKey string                 `db:"idempotency_key" json:"idempotency_key"`
    Payload        map[string]any         `db:"payload" json:"payload"`
    Headers        map[string]any         `db:"headers" json:"headers"`
    CreatedAt      time.Time              `db:"created_at" json:"created_at"`
    PublishedAt    *time.Time             `db:"published_at" json:"published_at"`
}

