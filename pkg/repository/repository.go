package repository

import (
	"github.com/jmoiron/sqlx"
	"orderTracking"
)

type Authorization interface {
	CreateUser(user orderTracking.User) (int, error)
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
