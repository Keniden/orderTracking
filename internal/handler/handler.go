package handler

import (
    "net/http"
    "orderTracking/internal/service"

    _ "orderTracking/docs"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/redis/go-redis/v9"
    otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Handler struct {
    services *service.Service
    rdb      *redis.Client
}

func NewHandler(services *service.Service, rdb *redis.Client) *Handler {
    return &Handler{services: services, rdb: rdb}
}

func (h *Handler) InitRoutes() *gin.Engine {
    router := gin.New()
    router.Use(otelgin.Middleware("order-tracking-api"))
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    router.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
    router.GET("/readyz", func(c *gin.Context) { c.String(http.StatusOK, "ready") })
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))

    auth := router.Group("/auth")
    {
        auth.POST("/sign-up", h.signUp)
        auth.POST("/sign-in", h.signIn)
    }
    api := router.Group("/api")
    {
        lists := api.Group("/lists", h.userIdentity, rateLimiter(h.rdb, 100, 60_000_000_000)) // 100 req/min
        {
            lists.POST("/add_order", h.createOrder)
            lists.GET("/status", h.getStatus)
        }
    }
    return router
}
