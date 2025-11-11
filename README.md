Order Tracking — Event-Driven Микросервисная Платформа

Ключевые технологии: Go, gRPC, grpc-gateway, Kafka, Redis, PostgreSQL, Neo4j, ClickHouse, Docker, Kubernetes, Prometheus, Grafana, OpenTelemetry, GitLab CI/CD

Архитектура
- Event-driven: доменные события публикуются в Kafka через outbox-паттерн (таблица outbox_events + воркер outbox-relay)
- Микросервисы: HTTP API (Gin), outbox-relay (публикация событий в Kafka), order-consumer (идемпотентный консьюмер, проекция в Redis/ClickHouse/Neo4j)
- Внутренние контракты: gRPC + grpc-gateway (прото-файлы в proto/, сервер с build-тегом with_grpc). Внешний REST уже доступен (Gin). 
- Хранения: PostgreSQL (источник истины, CQRS+outbox), Redis (кэш горячих агрегатов и rate limiting), ClickHouse (аналитика), Neo4j (граф доменных зависимостей)
- Наблюдаемость: Prometheus метрики (/metrics), OpenTelemetry (OTLP экспорт трейсов), готовые Docker/K8s манифесты

Компоненты
- cmd/main.go — HTTP API (создание заказа, получение статуса). Метрики Prometheus, трейсинг OTel, rate limiting через Redis.
- cmd/outbox-relay — воркер, публикующий события из Postgres outbox в Kafka (idempotency_key).
- cmd/order-consumer — Kafka-консьюмер с идемпотентной обработкой (Redis SETNX), пишет события в ClickHouse, кэширует статус заказа в Redis, проецирует узлы/связи в Neo4j.
- internal/kafka — минимальный Kafka producer (segmentio/kafka-go).
- internal/analytics — ClickHouse клиент и схема order_events.
- internal/graph — Neo4j клиент и базовая проекция графа.
- internal/telemetry — инициализация OpenTelemetry (OTLP gRPC).
- internal/metrics — счетчики Prometheus.
- internal/handler — REST ручки, middleware (аутентификация, rate limiting), /metrics, /healthz, /readyz.
- schema — миграции Postgres (users, orders, outbox_events).

Требования
- Go 1.25.0+
- Docker / Docker Compose (для локальной платформы)

Запуск локально (Docker Compose)
- Требования: Docker + Docker Compose
- Команда: docker compose -f deploy/docker-compose.yml up --build
- API: http://localhost:8000 (Swagger по /swagger/index.html), метрики /metrics
- Prometheus: http://localhost:9090, Grafana: http://localhost:3000 (анонимный доступ)
- Kafka: localhost:9092, PostgreSQL: localhost:5436, Redis: localhost:6379, ClickHouse: localhost:9000, Neo4j: localhost:7687

Сборка вручную
- go build ./cmd/...  (основные бинарники: api, outbox-relay, order-consumer)

gRPC и grpc-gateway
- Прото: proto/order.proto (HTTP аннотации для grpc-gateway)
- Генерация и сборка:
  - make proto   # генерирует pb и gateway-стабы
  - make grpc    # собирает gRPC сервер с тэгом with_grpc (порт gRPC 9090, HTTP gateway 8080)
- В CI (.gitlab-ci.yml) есть job build-grpc, который генерирует код и собирает gRPC-бинарник.

Kubernetes
- Манифесты в k8s/: Deployment/Service для API, Outbox Relay и Consumer, ConfigMap, HPA.

Тестирование
- Контрактные REST тесты: tests/rest_contract_test.go (httptest)
- Интеграционные: integration/outbox_test.go (Testcontainers; build tag integration)
  - go test -tags=integration ./integration -run TestOutboxInsert
- Нагрузка: load/order_create.js (k6)

Безопасность и качество
- GitLab CI/CD: .gitlab-ci.yml — lint, gosec, staticcheck, test, build, build-grpc

Важно
- Проект собирается и запускается как REST + event-driven пайплайн уже сейчас. gRPC/grpc-gateway включены на уровне контрактов (proto) и шаблонов и могут быть активированы генерацией кода (build-тег with_grpc не участвует в дефолтной сборке).
