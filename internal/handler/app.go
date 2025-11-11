package handler

import (
    "orderTracking/internal/repository"
    "orderTracking/internal/service"

    "github.com/jmoiron/sqlx"
    "github.com/redis/go-redis/v9"
)

func InitApp(db *sqlx.DB, rdb *redis.Client) *Handler {
    repos := repository.NewRepository(db)
    services := service.NewService(repos)
    handlers := NewHandler(services, rdb)
    return handlers
}
