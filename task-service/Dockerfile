FROM golang:1.25.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY task-service ./task-service
COPY notification-service ./notification-service
COPY proto ./proto

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /task-service ./task-service/cmd/task-service

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /task-service /task-service

EXPOSE 8080 50051

CMD ["/task-service"]