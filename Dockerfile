FROM golang:1.22.5-alpine AS builder

COPY . /github.com/nqxcode/auth_microservice/source/
WORKDIR /github.com/nqxcode/auth_microservice/source/

RUN go mod download
RUN go build -o ./bin/auth_microservice cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/nqxcode/auth_microservice/source/bin/auth_microservice .

CMD ["./auth_microservice"]