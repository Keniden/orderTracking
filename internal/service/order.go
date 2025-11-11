package service

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "orderTracking/internal/models"
    "orderTracking/internal/repository"
)

type OrderService struct {
    repo   repository.OrderItem
    outbox repository.Outbox
}

func NewOrderService(repo repository.OrderItem, outbox repository.Outbox) *OrderService {
    return &OrderService{repo: repo, outbox: outbox}
}

func (s *OrderService) CreateOrderItem(orderItem models.OrderItem) (string, error) {
    trackID, err := s.repo.CreateOrderItem(orderItem)
    if err != nil {
        return "", err
    }
    // Prepare outbox event for OrderCreated
    idemp := sha256.Sum256([]byte("order-created:" + trackID))
    _ = s.outbox.Enqueue(context.Background(), models.OutboxEvent{
        AggregateType:  "order",
        AggregateID:    trackID,
        EventType:      "OrderCreated",
        IdempotencyKey: hex.EncodeToString(idemp[:]),
        Payload: map[string]any{
            "track_id": trackID,
            "title":    orderItem.Title,
            "price":    orderItem.Price,
            "status":   orderItem.Status,
        },
    })
    return trackID, nil
}

func (s *OrderService) GetOrderDetailsByTrackID(track_id string) (models.OrderItem, error) {
    return s.repo.GetOrderDetailsByTrackID(track_id)
}
