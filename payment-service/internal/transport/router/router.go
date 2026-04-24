package router

import (
	"net/http"

	http_handlers "github.com/Xanaduxan/tasks-golang/payment-service/internal/transport/http-handlers"
	"github.com/Xanaduxan/tasks-golang/payment-service/internal/transport/http-handlers/middleware"
)

func New(jwtSecret []byte) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(
		"POST /payments",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.CreatePayment)),
	)
	mux.Handle(
		"GET /payments/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetPayment)),
	)

	mux.Handle(
		"PUT /payments/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.UpdatePayment)),
	)
	mux.Handle(
		"DELETE /payments/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.DeletePayment)),
	)
	mux.Handle(
		"POST /payments/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.ClosePayment)),
	)
	mux.Handle(
		"GET /shops/{shop_id}/payments",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetShopPayments)),
	)

	mux.Handle(
		"GET /shops",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetShops)),
	)
	mux.Handle(
		"GET /shop/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.GetShop)),
	)
	mux.Handle(
		"POST /shop",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.CreateShop)),
	)
	mux.Handle(
		"PUT /shop/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.UpdateShop)),
	)
	mux.Handle(
		"DELETE /shop/{id}",
		middleware.JWT(jwtSecret)(http.HandlerFunc(http_handlers.DeleteShop)),
	)

	return mux
}
