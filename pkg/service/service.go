package service

import "orderTracking/pkg/repository"

type Authorization interface {
}

type OrderTracking interface {
}

type OrderItem interface {
}

type Service struct {
	Authorization
	OrderTracking
	OrderItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
