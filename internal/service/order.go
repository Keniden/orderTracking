package service

import (
	"orderTracking/internal/models"
	"orderTracking/internal/repository"
)

type OrderService struct {
	repo repository.OrderItem
}

func NewOrderService(repo repository.OrderItem) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrderItem(orderItem models.OrderItem) (string, error) {
	return s.repo.CreateOrderItem(orderItem)
}

func (s *OrderService) GetOrderDetailsByTrackID(track_id string) (models.OrderItem, error) {
	return s.repo.GetOrderDetailsByTrackID(track_id)
}
