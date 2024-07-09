package models

import "time"

type OrderItem struct {
	Track_id    string    `json:"track_id"`
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Username    string    `json:"username"`
	Date        time.Time `json:"data"`
}

