package main

import (
    "context"
    "orderTracking/internal/api"
    "orderTracking/internal/cache"
    "orderTracking/internal/configs"
    "orderTracking/internal/handler"
    "orderTracking/internal/metrics"
    "orderTracking/internal/telemetry"

    _ "github.com/lib/pq"
    "github.com/prometheus/client_golang/prometheus"
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
    appCfg := configs.Load()
    // telemetry
    shutdown, err := telemetry.Setup("order-tracking-api", appCfg.Telemetry.OTLPEndpoint)
    if err != nil { logrus.Warnf("otel setup failed: %v", err) }
    defer func(){ if shutdown != nil { _ = shutdown(context.Background()) } }()

    db := configs.InitConfig()
    rdb := cache.NewRedis(appCfg.Redis.Addr, appCfg.Redis.Password, appCfg.Redis.DB)

    // metrics registry
    // register metrics on the default Prometheus registry to expose via /metrics
    _ = prometheus.NewRegistry() // keep a local registry if needed elsewhere
    metrics.RegisterDefault()

    handlers := handler.InitApp(db, rdb)
    server := new(api.Server)
    if err := server.Run(appCfg.Port, handlers.InitRoutes()); err != nil {
        logrus.Fatalf("error occured while running http server: %s", err.Error())
    }
}
