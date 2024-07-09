package main

import (
	"orderTracking/internal/api"
	"orderTracking/internal/configs"
	"orderTracking/internal/handler"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

//	@title			Order Tracking API
//	@version		1.0
//	@description	API Server OrderTracking Application

//@host localhost:8000
//@BasePath /

//@securityDefinitions.apikey ApiKeyAuth
//@in Header
//@name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	db := configs.InitConfig()
	handlers := handler.InitApp(db)
	server := new(api.Server)
	if err := server.Run("8000", handlers.InitRoutes()); err != nil {
		logrus.Fatal("error occured while running http server: %s", err.Error())
	}
}
