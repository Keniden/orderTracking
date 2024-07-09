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
	query := fmt.Sprintf("INSERT INTO %s (track_id, id, title, description, price, date, from_location, to_location, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING track_id", orderList)
	row := r.db.QueryRow(query, orderItem.Track_id, orderItem.Id, orderItem.Title, orderItem.Description, orderItem.Price, orderItem.Date, orderItem.From_location, orderItem.To_location, orderItem.Status)
	if err := row.Scan(&track_id); err != nil {
		return "", err
	}
	return track_id, nil
}

func (r *OrderRepository) GetOrderDetailsByTrackID(track_id string) (models.OrderItem, error) {
	var order models.OrderItem
	query := `SELECT track_id, id, title, description, price, date, from_location, to_location, status FROM orders WHERE track_id = $1`
	err := r.db.QueryRow(query, track_id).Scan(&order.Track_id, &order.Id, &order.Title, &order.Description, &order.Price, &order.Date, &order.From_location, &order.To_location, &order.Status)
	if err != nil {
		return order, err
	}
	return order, nil
}
