package repository

import (
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user order.User) (int, error)
}

type OrderList interface {
}

type OrderItem interface {
}

type Repository struct {
	Authorization
	OrderList
	OrderItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
