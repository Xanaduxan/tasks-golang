FROM golang:1.25.7 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/etl-worker/cmd/etl-worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /etl-worker

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /etl-worker /etl-worker

CMD ["/etl-worker"]
