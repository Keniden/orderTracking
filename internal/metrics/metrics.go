package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    OrdersCreated = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "orders", Subsystem: "api", Name: "created_total", Help: "Total created orders",
    })
    OrderStatusQueried = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "orders", Subsystem: "api", Name: "status_queries_total", Help: "Total status queries",
    })
    ConsumerProcessed = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "orders", Subsystem: "consumer", Name: "processed_total", Help: "Total events processed",
    })
    ConsumerSkipped = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "orders", Subsystem: "consumer", Name: "skipped_total", Help: "Total duplicate events skipped",
    })
    ConsumerErrors = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "orders", Subsystem: "consumer", Name: "errors_total", Help: "Total consumer errors",
    })
)

func MustRegister(reg *prometheus.Registry) {
    reg.MustRegister(OrdersCreated, OrderStatusQueried, ConsumerProcessed, ConsumerSkipped, ConsumerErrors)
}

func RegisterDefault() {
    prometheus.MustRegister(OrdersCreated, OrderStatusQueried, ConsumerProcessed, ConsumerSkipped, ConsumerErrors)
}
