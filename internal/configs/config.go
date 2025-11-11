package configs

import (
    "log"
    "orderTracking/internal/repository"

    "github.com/jmoiron/sqlx"
    "github.com/sirupsen/logrus"
    "github.com/spf13/viper"
)

type KafkaConfig struct {
    Brokers []string
    OrderTopic string
}

type RedisConfig struct {
    Addr string
    DB   int
    Password string
}

type ClickHouseConfig struct {
    Addr string
    Database string
}

type Neo4jConfig struct {
    URI      string
    Username string
    Password string
}

type TelemetryConfig struct {
    OTLPEndpoint string
}

type AppConfig struct {
    Port string
    Kafka KafkaConfig
    Redis RedisConfig
    ClickHouse ClickHouseConfig
    Neo4j Neo4jConfig
    Telemetry TelemetryConfig
}

func initViper() error {
    viper.AddConfigPath("internal/configs")
    viper.SetConfigName("config")
    return viper.ReadInConfig()
}

func InitConfig() *sqlx.DB {
    if err := initViper(); err != nil {
        logrus.Fatalf("error occured while initializing config:  %s", err.Error())
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

func Load() AppConfig {
    if err := initViper(); err != nil {
        logrus.Fatalf("error occured while initializing config:  %s", err.Error())
    }
    return AppConfig{
        Port: viper.GetString("port"),
        Kafka: KafkaConfig{
            Brokers: viper.GetStringSlice("kafka.brokers"),
            OrderTopic: viper.GetString("kafka.order_topic"),
        },
        Redis: RedisConfig{
            Addr: viper.GetString("redis.addr"),
            DB: viper.GetInt("redis.db"),
            Password: viper.GetString("redis.password"),
        },
        ClickHouse: ClickHouseConfig{
            Addr: viper.GetString("clickhouse.addr"),
            Database: viper.GetString("clickhouse.database"),
        },
        Neo4j: Neo4jConfig{
            URI: viper.GetString("neo4j.uri"),
            Username: viper.GetString("neo4j.username"),
            Password: viper.GetString("neo4j.password"),
        },
        Telemetry: TelemetryConfig{
            OTLPEndpoint: viper.GetString("telemetry.otlp_endpoint"),
        },
    }
}
