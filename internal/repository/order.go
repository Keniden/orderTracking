package repository

import (
	"fmt"
	"orderTracking/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}
func (r *OrderRepository) CreateOrderItem(orderItem models.OrderItem) (string, error) {
	var track_id string
	query := fmt.Sprintf("INSERT INTO %s (track_id, id, title, description, price, username, date) values ($1, $2, $3, $4, $5, $6, $7) RETURNING track_id", orderList)
	row := r.db.QueryRow(query, orderItem.Track_id, orderItem.Id, orderItem.Title, orderItem.Description,orderItem.Price, orderItem.Username, orderItem.Date)
	if err := row.Scan(&track_id); err != nil {
		return "", err
	}
	return track_id, nil
}
