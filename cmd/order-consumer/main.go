package main

import (
    "context"
    "encoding/json"
    "log"
    "strings"
    "time"

    "orderTracking/internal/analytics"
    "orderTracking/internal/cache"
    "orderTracking/internal/configs"
    "orderTracking/internal/graph"
    "orderTracking/internal/metrics"

    ch "github.com/ClickHouse/clickhouse-go/v2"
    kfk "github.com/segmentio/kafka-go"
)

func main() {
    cfg := configs.Load()
    ctx := context.Background()

    // Metrics are process-wide; nothing to expose here but counters will be scraped if added to a server
    _ = ch.Settings{} // silence unused import in some environments

    r := kfk.NewReader(kfk.ReaderConfig{
        Brokers: cfg.Kafka.Brokers,
        GroupID: "orders-consumer",
        Topic:   cfg.Kafka.OrderTopic,
        MinBytes: 1,
        MaxBytes: 10e6,
    })
    defer r.Close()

    rdb := cache.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

    click, err := analytics.Connect(cfg.ClickHouse.Addr, cfg.ClickHouse.Database)
    if err != nil { log.Fatalf("clickhouse connect: %v", err) }
    if err := click.EnsureSchema(ctx); err != nil { log.Fatalf("clickhouse schema: %v", err) }

    g, err := graph.Connect(cfg.Neo4j.URI, cfg.Neo4j.Username, cfg.Neo4j.Password)
    if err != nil { log.Printf("neo4j connect failed: %v", err) } // optional
    defer func(){ if g != nil { _ = g.Close(ctx) } }()

    for {
        m, err := r.FetchMessage(ctx)
        if err != nil { log.Printf("fetch: %v", err); time.Sleep(time.Second); continue }
        idemp := header(&m, "idempotency_key")
        if idemp == "" { idemp = string(m.Key) + ":" + string(m.Value) }

        ok, err := cache.SetNX(ctx, rdb, "idemp:"+idemp, "1", 24*time.Hour)
        if err != nil { log.Printf("redis error: %v", err); metrics.ConsumerErrors.Inc(); continue }
        if !ok { metrics.ConsumerSkipped.Inc(); _ = r.CommitMessages(ctx, m); continue }

        var payload map[string]any
        _ = json.Unmarshal(m.Value, &payload)
        aggType := str(payload["aggregate_type"]) 
        aggID := str(payload["aggregate_id"]) 
        evType := str(payload["event_type"]) 
        if b, err := json.Marshal(payload); err == nil {
            if err := click.InsertEvent(ctx, idemp, aggType, aggID, evType, string(b)); err != nil {
                log.Printf("click insert: %v", err); metrics.ConsumerErrors.Inc(); continue
            }
        }

        // Update hot caches
        if status := str(payload["status"]); status != "" && aggID != "" {
            rdb.Set(ctx, "order:status:"+aggID, status, 24*time.Hour)
            rdb.Set(ctx, "order:last_event:"+aggID, string(m.Value), 24*time.Hour)
        }

        // Project basic path into graph if available
        if g != nil {
            from := str(payload["from_location"]) ; to := str(payload["to_location"]) ;
            if from != "" && to != "" && aggID != "" { _ = g.UpsertOrderPath(ctx, aggID, from, to) }
        }

        if err := r.CommitMessages(ctx, m); err != nil { log.Printf("commit: %v", err) }
        metrics.ConsumerProcessed.Inc()
    }
}

func header(m *kfk.Message, key string) string {
    for _, h := range m.Headers { if strings.EqualFold(h.Key, key) { return string(h.Value) } }
    return ""
}

func str(v any) string { if v == nil { return "" }; if s, ok := v.(string); ok { return s }; return "" }

