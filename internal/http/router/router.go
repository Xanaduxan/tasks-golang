package router

import (
	"net/http"

	"github.com/Xanaduxan/tasks-golang/internal/http/handlers"
	"github.com/Xanaduxan/tasks-golang/internal/http/handlers/middleware"
)

func New(jwtSecret []byte) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", handlers.Login)
	mux.HandleFunc("POST /registration", handlers.Registration)

	mux.Handle(
		"GET /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.GetTask)),
	)
	mux.Handle(
		"POST /task",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.CreateTask)),
	)
	mux.Handle(
		"PUT /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.UpdateTask)),
	)
	mux.Handle(
		"DELETE /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.DeleteTask)),
	)
	return mux
}
