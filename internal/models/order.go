package models

import "time"

type OrderItem struct {
	Track_id      string    `json:"track_id"`
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         int       `json:"price"`
	Date          time.Time `json:"data"`
	From_location string    `json:"from_location"`
	To_location   string    `json:"to_location"`
	Status        string    `json:"status"`
}
