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

	mux.Handle(
		"GET /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.GetProduct)),
	)
	mux.Handle(
		"POST /product",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.CreateProduct)),
	)
	mux.Handle(
		"PUT /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.UpdateProduct)),
	)
	mux.Handle(
		"DELETE /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.DeleteProduct)),
	)

	mux.Handle(
		"GET /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.GetDelivery)),
	)
	mux.Handle(
		"POST /delivery",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.CreateDelivery)),
	)
	mux.Handle(
		"PUT /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.UpdateDelivery)),
	)
	mux.Handle(
		"DELETE /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.DeleteDelivery)),
	)
	mux.Handle(
		"GET /stock",
		middleware.JWT(jwtSecret)(http.HandlerFunc(handlers.GetStocks)),
	)

	return mux
}
