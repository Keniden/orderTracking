package service

import (
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

type Service struct {
	Authorization
	OrderItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		OrderItem:     NewOrderService(repos.OrderItem),
	}
}
