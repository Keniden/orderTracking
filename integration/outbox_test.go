//go:build integration

package integration

import (
    "context"
    "os"
    "path/filepath"
    "testing"

    "orderTracking/internal/models"
    "orderTracking/internal/repository"

    tc "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/jmoiron/sqlx"
)

func TestOutboxInsert(t *testing.T) {
    ctx := context.Background()
    pg, err := postgres.RunContainer(ctx,
        postgres.WithDatabase("postgres"),
        postgres.WithUsername("postgres"),
        postgres.WithPassword("qwerty"),
    )
    if err != nil { t.Fatalf("postgres: %v", err) }
    defer pg.Terminate(ctx)

    host, _ := pg.Host(ctx)
    port, _ := pg.MappedPort(ctx, "5432")

    db, err := repository.NewPostgresDB(repository.Config{
        Host: host,
        Port: port.Port(),
        Username: "postgres",
        Password: "qwerty",
        DBName: "postgres",
        SSLMode: "disable",
    })
    if err != nil { t.Fatalf("db: %v", err) }

    // Apply minimal migrations: outbox schema
    applySQL(t, db, mustRead(t, filepath.Join("schema", "000002_outbox.up.sql")))

    outbox := repository.NewOutboxRepository(db)
    err = outbox.Enqueue(ctx, models.OutboxEvent{
        AggregateType: "order", AggregateID: "tid", EventType: "OrderCreated",
        IdempotencyKey: "k", Payload: map[string]any{"track_id":"tid"},
    })
    if err != nil { t.Fatalf("enqueue: %v", err) }
}

func applySQL(t *testing.T, db *sqlx.DB, sql string) {
    t.Helper()
    if _, err := db.ExecContext(context.Background(), sql); err != nil { t.Fatalf("migrate: %v", err) }
}

func mustRead(t *testing.T, p string) string {
    t.Helper()
    b, err := os.ReadFile(p)
    if err != nil { t.Fatalf("read: %v", err) }
    return string(b)
}
