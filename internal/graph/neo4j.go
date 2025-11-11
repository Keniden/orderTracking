package graph

import (
    "context"

    neo4j "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Graph struct {
    drv neo4j.DriverWithContext
}

func Connect(uri, user, pass string) (*Graph, error) {
    drv, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
    if err != nil { return nil, err }
    return &Graph{drv: drv}, nil
}

func (g *Graph) Close(ctx context.Context) error { return g.drv.Close(ctx) }

// UpsertOrderPath ensures a minimal graph: (Order)-[:ROUTES_TO]->(Hub)
func (g *Graph) UpsertOrderPath(ctx context.Context, trackID, from, to string) error {
    sess := g.drv.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer sess.Close(ctx)
    _, err := sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        _, err := tx.Run(ctx, `
            MERGE (o:Order {track_id: $track_id})
            SET o.updated_at = timestamp()
            MERGE (a:Hub {name: $from})
            MERGE (b:Hub {name: $to})
            MERGE (o)-[:ROUTES_FROM]->(a)
            MERGE (o)-[:ROUTES_TO]->(b)
        `, map[string]any{"track_id": trackID, "from": from, "to": to})
        return nil, err
    })
    return err
}

type Edge struct {
    From  string
    To    string
    Count int64
}

// CriticalEdgesByTraffic returns edges between hubs with order count >= minCount
func (g *Graph) CriticalEdgesByTraffic(ctx context.Context, minCount int) ([]Edge, error) {
    sess := g.drv.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer sess.Close(ctx)
    res, err := sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
        q := `
            MATCH (o:Order)-[:ROUTES_FROM]->(a:Hub)
            MATCH (o)-[:ROUTES_TO]->(b:Hub)
            WITH a,b, count(o) AS cnt
            WHERE cnt >= $min
            RETURN a.name AS from, b.name AS to, cnt AS cnt
            ORDER BY cnt DESC
        `
        rows, err := tx.Run(ctx, q, map[string]any{"min": minCount})
        if err != nil { return nil, err }
        var out []Edge
        for rows.Next(ctx) {
            rec := rows.Record()
            from, _ := rec.Get("from")
            to, _ := rec.Get("to")
            cnt, _ := rec.Get("cnt")
            out = append(out, Edge{From: from.(string), To: to.(string), Count: cnt.(int64)})
        }
        return out, rows.Err()
    })
    if err != nil { return nil, err }
    return res.([]Edge), nil
}
