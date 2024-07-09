package configs

import (
	"log"
	"orderTracking/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initViper() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func InitConfig() *sqlx.DB {
	if err := initViper(); err != nil {
		logrus.Fatal("error occured while initializing config:  %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
	})
	if err != nil {
		log.Fatalf("Error initialization db: %s", err.Error())
	}
	return db
}
