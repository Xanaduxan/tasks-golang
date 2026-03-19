package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TasksCurrent = prometheus.NewGauge(prometheus.GaugeOpts{Namespace: "app",
		Subsystem: "tasks",
		Name:      "current",
		Help:      "Current number of tasks."})
	TaskProcessingDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "app",
		Subsystem: "tasks",
		Name:      "processing_duration_seconds",
		Help:      "Task processing duration in seconds.",
		Buckets:   []float64{0.1, 0.3, 0.5, 1, 2, 5, 10, 30, 60},
	})
)

func init() {
	prometheus.MustRegister(TasksCurrent, TaskProcessingDuration)
}
