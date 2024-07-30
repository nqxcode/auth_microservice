ENV_FILE ?= .env

include $(ENV_FILE)

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(POSTGRES_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-auth-api

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/auth_v1/auth.proto

build:
	make build-grpc-server

build-grpc-server:
	make build-target TARGET=grpc_server

build-target:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o ./bin/$(TARGET) ./cmd/$(TARGET)

copy-to-server:
	scp ./bin/grpc_server root@${REGISTRY_HOST}:/root/auth_grpc_server

install-docker-buildx:
	mkdir -p ~/.docker/cli-plugins && \
    curl -L https://github.com/docker/buildx/releases/download/v0.16.0/buildx-v0.16.0.linux-amd64 -o ~/.docker/cli-plugins/docker-buildx && \
    chmod +x ~/.docker/cli-plugins/docker-buildx

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/nqxcode/auth-microservice:v0.0.2 .
	docker login -u token -p ${REGISTRY_PASSWORD} cr.selcloud.ru/nqxcode
	docker push cr.selcloud.ru/nqxcode/auth-microservice:v0.0.1


local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v
