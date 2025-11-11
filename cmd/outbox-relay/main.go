package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "orderTracking/internal/configs"
    kprod "orderTracking/internal/kafka"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func main() {
    // load config and DB
    cfg := configs.Load()
    db := configs.InitConfig()

    producer := kprod.NewProducer(kprod.Config{Brokers: cfg.Kafka.Brokers, Topic: cfg.Kafka.OrderTopic})
    defer producer.Close()

    ctx := context.Background()
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        if err := publishBatch(ctx, db, producer); err != nil {
            log.Printf("outbox publish error: %v", err)
        }
    }
}

type outboxRow struct {
    ID             int64          `db:"id"`
    AggregateType  string         `db:"aggregate_type"`
    AggregateID    string         `db:"aggregate_id"`
    EventType      string         `db:"event_type"`
    IdempotencyKey string         `db:"idempotency_key"`
    Payload        json.RawMessage `db:"payload"`
    Headers        json.RawMessage `db:"headers"`
}

func publishBatch(ctx context.Context, db *sqlx.DB, p *kprod.Producer) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    rows := []outboxRow{}
    err = tx.SelectContext(ctx, &rows, `SELECT id, aggregate_type, aggregate_id, event_type, idempotency_key, payload, headers
        FROM outbox_events
        WHERE published_at IS NULL
        ORDER BY id
        LIMIT 50
        FOR UPDATE SKIP LOCKED`)
    if err != nil { return err }
    for _, r := range rows {
        msg := map[string]any{
            "aggregate_type": r.AggregateType,
            "aggregate_id":   r.AggregateID,
            "event_type":     r.EventType,
        }
        var payload map[string]any
        _ = json.Unmarshal(r.Payload, &payload)
        for k, v := range payload { msg[k] = v }
        b, _ := json.Marshal(msg)
        headers := map[string]string{"idempotency_key": r.IdempotencyKey}
        if err := p.Publish(ctx, []byte(r.AggregateID), b, headers); err != nil {
            return err
        }
        if _, err := tx.ExecContext(ctx, `UPDATE outbox_events SET published_at = now() WHERE id=$1`, r.ID); err != nil {
            return err
        }
    }
    return tx.Commit()
}

