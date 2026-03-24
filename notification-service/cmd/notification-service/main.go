package main

import (
	"log"
	"net"
	"net/http"

	grpcserver "github.com/Xanaduxan/tasks-golang/notification-service/internal/transport/grpc"
	"github.com/Xanaduxan/tasks-golang/notification-service/internal/transport/websocket"
	notificationpb "github.com/Xanaduxan/tasks-golang/notification-service/pkg/pb/notification/v1"
	"google.golang.org/grpc"
)

func main() {
	wsManager := websocket.NewManager()
	wsHandler := websocket.NewHandler(wsManager)

	// gRPC
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("failed to listen grpc: %v", err)
		}

		s := grpc.NewServer()
		notificationpb.RegisterNotificationServiceServer(
			s,
			grpcserver.NewServer(wsManager),
		)

		log.Println("notification gRPC listening on :50052")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("grpc error: %v", err)
		}
	}()

	// HTTP (websocket)
	mux := http.NewServeMux()
	mux.Handle("GET /ws", wsHandler)

	log.Println("notification http listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
