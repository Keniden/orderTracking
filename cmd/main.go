package main

import (
	"context"
	"net/http"
	"orderTracking/internal/configs"
	"orderTracking/internal/handler"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 mb
		ReadTimeout:    10 * time.Second,
		IdleTimeout:    10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	db := configs.InitConfig()
	handlers := handler.InitApp(db)
	server := new(Server)
	if err := server.Run("8000", handlers.InitRoutes()); err != nil {
		logrus.Fatal("error occured while running http server: %s", err.Error())
	}
}
