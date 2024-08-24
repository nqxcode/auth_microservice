ENV_FILE ?= .env

include $(ENV_FILE)

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(POSTGRES_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

GIT_SHA=$(shell git rev-parse --short HEAD)

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.1.3

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

login-to-server:
	@ssh root@${SERVER_HOST}

copy-to-server:
	scp ./bin/grpc_server root@${REGISTRY}:/root/auth_grpc_server

show-git-sha:
	@echo $(GIT_SHA)

install-docker-buildx:
	mkdir -p ~/.docker/cli-plugins && \
    curl -L https://github.com/docker/buildx/releases/download/v0.16.0/buildx-v0.16.0.linux-amd64 -o ~/.docker/cli-plugins/docker-buildx && \
    chmod +x ~/.docker/cli-plugins/docker-buildx

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/nqxcode/auth-microservice:$(GIT_SHA) -f ./Dockerfile .
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/nqxcode/auth-microservice-migration:$(GIT_SHA) -f ./migration.Dockerfile .
	docker login -u token -p ${REGISTRY_PASSWORD} $(REGISTRY)
	docker push $(REGISTRY)/auth-microservice:$(GIT_SHA)
	docker push $(REGISTRY)/auth-microservice-migration:$(GIT_SHA)


local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-create:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} create "$(call args)" sql

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

.PHONY: test
test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/nqxcode/auth_microservice/internal/service/...,github.com/nqxcode/auth_microservice/internal/api/... -count 5

.PHONY: test-coverage
test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/nqxcode/auth_microservice/internal/service/...,github.com/nqxcode/auth_microservice/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore