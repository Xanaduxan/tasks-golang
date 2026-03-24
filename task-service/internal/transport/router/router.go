package router

import (
	"net/http"

	http_handlers2 "github.com/Xanaduxan/tasks-golang/task-service/internal/transport/http-handlers"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/transport/http-handlers/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(jwtSecret []byte) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", http_handlers2.Login)
	mux.HandleFunc("POST /registration", http_handlers2.Registration)

	mux.Handle(
		"GET /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.GetTask)),
	)
	mux.Handle(
		"POST /task",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.CreateTask)),
	)
	mux.Handle(
		"GET /task/search",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.SearchTasks)),
	)
	mux.Handle(
		"PUT /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.UpdateTask)),
	)
	mux.Handle(
		"DELETE /task/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.DeleteTask)),
	)

	mux.Handle(
		"GET /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.GetProduct)),
	)
	mux.Handle(
		"POST /product",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.CreateProduct)),
	)
	mux.Handle(
		"PUT /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.UpdateProduct)),
	)
	mux.Handle(
		"DELETE /product/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.DeleteProduct)),
	)

	mux.Handle(
		"GET /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.GetDelivery)),
	)
	mux.Handle(
		"POST /delivery",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.CreateDelivery)),
	)
	mux.Handle(
		"PUT /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.UpdateDelivery)),
	)
	mux.Handle(
		"DELETE /delivery/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.DeleteDelivery)),
	)

	mux.Handle(
		"GET /stock",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.GetStocks)),
	)

	mux.Handle(
		"GET /group/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.GetGroup)),
	)
	mux.Handle(
		"POST /group",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.CreateGroup)),
	)

	mux.Handle(
		"GET /group/{group_id}/member",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.GetGroupMembers)),
	)
	mux.Handle(
		"POST /group/{group_id}/member",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.CreateGroupMember)),
	)
	mux.Handle(
		"PUT /group/{group_id}/member/{user_id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.UpdateGroupMember)),
	)
	mux.Handle(
		"DELETE /group/{group_id}/member/{user_id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers2.DeleteGroupMember)),
	)

	//mux.Handle("GET /ws", wsHandler)
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
