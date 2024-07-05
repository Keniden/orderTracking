package main

import (
	"log"
	"orderTracking"
	"orderTracking/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)

	server := new(orderTracking.Server)
	if err := server.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatal("error occured while running http server: %s", err.Error())
	}
}
