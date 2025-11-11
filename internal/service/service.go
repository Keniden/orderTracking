package service

import (
    "context"
    "orderTracking/internal/models"
    "orderTracking/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type OrderItem interface {
    CreateOrderItem(orderItem models.OrderItem) (string, error)
    GetOrderDetailsByTrackID(track_id string) (models.OrderItem, error)
}

type Outbox interface {
    Enqueue(ctx context.Context, ev models.OutboxEvent) error
}

type Service struct {
    Authorization
    OrderItem
    Outbox
}

func NewService(repos *repository.Repository) *Service {
    return &Service{
        Authorization: NewAuthService(repos.Authorization),
        OrderItem:     NewOrderService(repos.OrderItem, repos.Outbox),
        Outbox:        repos.Outbox,
    }
}
