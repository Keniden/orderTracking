package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "orderTracking/internal/models"

    "github.com/jmoiron/sqlx"
)

type Outbox interface {
    Enqueue(ctx context.Context, ev models.OutboxEvent) error
}

type OutboxRepository struct {
    db *sqlx.DB
}

func NewOutboxRepository(db *sqlx.DB) *OutboxRepository {
    return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Enqueue(ctx context.Context, ev models.OutboxEvent) error {
    payloadJSON, err := json.Marshal(ev.Payload)
    if err != nil {
        return fmt.Errorf("marshal payload: %w", err)
    }
    headersJSON, err := json.Marshal(ev.Headers)
    if err != nil {
        return fmt.Errorf("marshal headers: %w", err)
    }
    q := `INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, idempotency_key, payload, headers)
          VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb)`
    _, err = r.db.ExecContext(ctx, q, ev.AggregateType, ev.AggregateID, ev.EventType, ev.IdempotencyKey, string(payloadJSON), string(headersJSON))
    return err
}

