package analytics

import (
    "context"
    "time"

    ch "github.com/ClickHouse/clickhouse-go/v2"
)

type Click struct {
    conn ch.Conn
}

func Connect(addr, database string) (*Click, error) {
    conn, err := ch.Open(&ch.Options{
        Addr: []string{addr},
        Auth: ch.Auth{Database: database},
        DialTimeout: 5 * time.Second,
        Debug: false,
        ClientInfo: ch.ClientInfo{Products: []struct{ Name, Version string }{{Name: "order-tracking", Version: "0.1"}}},
    })
    if err != nil { return nil, err }
    return &Click{conn: conn}, nil
}

func (c *Click) EnsureSchema(ctx context.Context) error {
    ddl := `CREATE TABLE IF NOT EXISTS order_events(
        idempotency_key String,
        aggregate_type String,
        aggregate_id String,
        event_type String,
        ts DateTime DEFAULT now(),
        payload String
    ) ENGINE = MergeTree ORDER BY (aggregate_type, aggregate_id, ts)`
    return c.conn.Exec(ctx, ddl)
}

func (c *Click) InsertEvent(ctx context.Context, idemp, aggType, aggID, evType, payload string) error {
    batch, err := c.conn.PrepareBatch(ctx, `INSERT INTO order_events (idempotency_key, aggregate_type, aggregate_id, event_type, payload) VALUES`)
    if err != nil { return err }
    if err := batch.Append(idemp, aggType, aggID, evType, payload); err != nil { return err }
    return batch.Send()
}

