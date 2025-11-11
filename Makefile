PROTO_DIR=proto
GEN_DIR=internal/gen

.PHONY: proto tools build api relay consumer grpc run-compose

tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

proto: tools
	protoc -I . --go_out=. --go-grpc_out=. --grpc-gateway_out=. $(PROTO_DIR)/order.proto

build:
	go build ./...

api:
	go build -o bin/api ./cmd

relay:
	go build -o bin/outbox-relay ./cmd/outbox-relay

consumer:
	go build -o bin/order-consumer ./cmd/order-consumer

grpc: proto
	go build -tags with_grpc -o bin/grpc-server ./cmd/grpc-server

run-compose:
	docker compose -f deploy/docker-compose.yml up --build

