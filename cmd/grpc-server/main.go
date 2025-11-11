//go:build with_grpc

package main

import (
    "context"
    "log"
    "net"
    "net/http"
    "orderTracking/internal/configs"
    "orderTracking/internal/repository"
    grpcserver "orderTracking/internal/server/grpc"
    "orderTracking/internal/service"

    orderv1 "orderTracking/internal/gen/ordertracking/v1"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    _ "github.com/lib/pq"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

func main() {
    // Load config and init dependencies
    cfg := configs.Load()
    db := configs.InitConfig()
    repos := repository.NewRepository(db)
    svc := service.NewService(repos)

    // Start gRPC server
    grpcLis, err := net.Listen("tcp", ":9090")
    if err != nil { log.Fatalf("listen: %v", err) }
    grpcServer := grpc.NewServer()
    srv := grpcserver.NewOrderServer(svc)
    orderv1.RegisterOrderServiceServer(grpcServer, srv)
    reflection.Register(grpcServer)
    go func(){ _ = grpcServer.Serve(grpcLis) }()

    // Start grpc-gateway HTTP server mapping to gRPC endpoint
    mux := runtime.NewServeMux()
    ctx := context.Background()
    err = orderv1.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:9090", []grpc.DialOption{grpc.WithInsecure()})
    if err != nil { log.Fatalf("gateway: %v", err) }
    log.Fatal(http.ListenAndServe(":8080", mux))
}
