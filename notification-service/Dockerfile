FROM golang:1.25.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY task-service ./task-service
COPY notification-service ./notification-service
COPY proto ./proto

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./notification-service/cmd/notification-service

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /notification-service /notification-service

EXPOSE 8081 50052

CMD ["/notification-service"]