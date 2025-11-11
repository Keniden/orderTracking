//go:build with_grpc

package grpcserver

import (
    "context"
    "orderTracking/internal/models"
    "orderTracking/internal/service"
    orderv1 "orderTracking/internal/gen/ordertracking/v1"
)

type OrderServer struct{
    svc *service.Service
    orderv1.UnimplementedOrderServiceServer
}

func NewOrderServer(svc *service.Service) *OrderServer { return &OrderServer{svc: svc} }

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
    // In gRPC path, user identity should come from auth metadata; for demo use id=1
    o := models.OrderItem{
        Title: req.GetTitle(), Description: req.GetDescription(), Price: int(req.GetPrice()), From_location: req.GetFromLocation(), To_location: req.GetToLocation(), Status: "created",
        Id: 1,
    }
    trackID, err := s.svc.CreateOrderItem(o)
    if err != nil { return nil, err }
    return &orderv1.CreateOrderResponse{TrackId: trackID}, nil
}

func (s *OrderServer) GetStatus(ctx context.Context, req *orderv1.GetStatusRequest) (*orderv1.GetStatusResponse, error) {
    item, err := s.svc.GetOrderDetailsByTrackID(req.GetTrackId())
    if err != nil { return nil, err }
    return &orderv1.GetStatusResponse{TrackId: item.Track_id, Status: item.Status}, nil
}

