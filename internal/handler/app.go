package handler

import (
	"orderTracking/internal/repository"
	"orderTracking/internal/service"

	"github.com/jmoiron/sqlx"
)

func InitApp(db *sqlx.DB) *Handler {
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := NewHandler(services)
	return handlers
}
