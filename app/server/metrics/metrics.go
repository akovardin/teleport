package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Rps = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "rps",
			Namespace: "progres",
			Subsystem: "requests",
			Help:      "Rps.",
		},
		[]string{"path", "success", "code"},
	)

	Timings = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "timings",
			Namespace:  "progres",
			Subsystem:  "requests",
			Help:       "Timings.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(Rps)
	prometheus.MustRegister(Timings)
}
