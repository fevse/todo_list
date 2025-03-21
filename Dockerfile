FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

# RUN apk add --no-cache \
#     # Important: required for go-sqlite3
#     gcc \
#     # Required for Alpine
#     musl-dev

COPY go.mod .
COPY . .

RUN go mod download
RUN go build -o todolist cmd/todolist/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/todolist .
COPY --from=builder /app/configs/config.toml .
COPY --from=builder /app/.env .
COPY --from=builder /app/db/migrations/20240627182704_todostor.sql .

CMD ["./todolist"]