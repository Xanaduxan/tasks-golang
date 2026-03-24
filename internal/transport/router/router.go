package router

import (
	"net/http"

	"github.com/Xanaduxan/tasks-golang/internal/transport/http-handlers"
	"github.com/Xanaduxan/tasks-golang/internal/transport/http-handlers/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(jwtSecret []byte, wsHandler http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", http_handlers.Login)
	mux.HandleFunc("POST /registration", http_handlers.Registration)

	mux.Handle(
		"GET /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetTask)),
	)
	mux.Handle(
		"POST /task",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.CreateTask)),
	)
	mux.Handle(
		"GET /task/search",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.SearchTasks)),
	)
	mux.Handle(
		"PUT /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.UpdateTask)),
	)
	mux.Handle(
		"DELETE /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.DeleteTask)),
	)

	mux.Handle(
		"GET /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetProduct)),
	)
	mux.Handle(
		"POST /product",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.CreateProduct)),
	)
	mux.Handle(
		"PUT /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.UpdateProduct)),
	)
	mux.Handle(
		"DELETE /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.DeleteProduct)),
	)

	mux.Handle(
		"GET /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetDelivery)),
	)
	mux.Handle(
		"POST /delivery",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.CreateDelivery)),
	)
	mux.Handle(
		"PUT /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.UpdateDelivery)),
	)
	mux.Handle(
		"DELETE /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.DeleteDelivery)),
	)

	mux.Handle(
		"GET /stock",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetStocks)),
	)

	mux.Handle(
		"GET /group/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetGroup)),
	)
	mux.Handle(
		"POST /group",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.CreateGroup)),
	)

	mux.Handle(
		"GET /group/{group_id}/member",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetGroupMembers)),
	)
	mux.Handle(
		"POST /group/{group_id}/member",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.CreateGroupMember)),
	)
	mux.Handle(
		"PUT /group/{group_id}/member/{user_id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.UpdateGroupMember)),
	)
	mux.Handle(
		"DELETE /group/{group_id}/member/{user_id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.DeleteGroupMember)),
	)

	mux.Handle("GET /ws", wsHandler)
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
