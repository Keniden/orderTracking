//go:build with_grpc

package tests

import (
    "context"
    "net"
    "testing"

    orderv1 "orderTracking/internal/gen/ordertracking/v1"
    grpcserver "orderTracking/internal/server/grpc"
    "orderTracking/internal/models"
    "orderTracking/internal/service"

    "google.golang.org/grpc"
)

type fakeAuth2 struct{ fakeAuth }
type fakeOrder2 struct{ fakeOrder }

func TestGRPCCreateOrder(t *testing.T) {
    svc := &service.Service{Authorization: fakeAuth2{}, OrderItem: fakeOrder2{}}
    srv := grpcserver.NewOrderServer(svc)

    lis, err := net.Listen("tcp", ":0")
    if err != nil { t.Fatal(err) }
    gs := grpc.NewServer()
    orderv1.RegisterOrderServiceServer(gs, srv)
    go gs.Serve(lis)
    defer gs.GracefulStop()

    conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
    if err != nil { t.Fatal(err) }
    defer conn.Close()
    c := orderv1.NewOrderServiceClient(conn)
    res, err := c.CreateOrder(context.Background(), &orderv1.CreateOrderRequest{Title: "t"})
    if err != nil { t.Fatalf("rpc error: %v", err) }
    if res.GetTrackId() == "" { t.Fatalf("empty track id") }
}

