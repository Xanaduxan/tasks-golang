FROM golang:1.25.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY payment-service ./payment-service
COPY notification-service ./notification-service
COPY proto ./proto

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /payment-service ./payment-service/cmd/payment-service

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /payment-service /payment-service

EXPOSE 8080 50051

CMD ["/payment-service"]