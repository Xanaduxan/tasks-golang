package grpc

import (
	"net"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/tasks"
	grpcHandlers "github.com/Xanaduxan/tasks-golang/task-service/internal/transport/grpc/handlers"
	taskv1 "github.com/Xanaduxan/tasks-golang/task-service/pkg/pb/task/v1"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *ggrpc.Server
	addr       string
}

func NewServer(addr string, taskService *tasks.Service) *Server {
	grpcServer := ggrpc.NewServer()

	taskHandler := grpcHandlers.NewTaskHandler(taskService)
	taskv1.RegisterTaskServiceServer(grpcServer, taskHandler)
	reflection.Register(grpcServer)
	return &Server{
		grpcServer: grpcServer,
		addr:       addr,
	}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}
