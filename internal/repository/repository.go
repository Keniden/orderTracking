package repository

import (
	"orderTracking/internal/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type OrderItem interface {
	CreateOrderItem(orderItem models.OrderItem) (string, error)
}

type Repository struct {
	Authorization
	OrderItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		OrderItem:     NewOrderRepository(db),
	}
}
