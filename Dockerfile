FROM golang:1.25.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env /app/.env

EXPOSE 8080

CMD ["/app/server"]