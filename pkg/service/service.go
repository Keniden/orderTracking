package service

import (
	"orderTracking/pkg/repository"
	"orderTracking"
)

type Authorization interface {
	CreateUser(user orderTracking.User) (int, error)
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
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
