package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	WebSocketConnectionsActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "app",
		Subsystem: "websocket",
		Name:      "connections_active",
		Help:      "Current number of active WebSocket connections.",
	})
)

func init() {
	log.Println("registering websocket metric")
	prometheus.MustRegister(WebSocketConnectionsActive)
}
