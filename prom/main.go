package prom

import "github.com/prometheus/client_golang/prometheus"

var HomeHandlerCounter prometheus.Counter
var PathLength *prometheus.HistogramVec

func Init() {

	HomeHandlerCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "",
		Name:      "Home_handler_counter",
		Help:      "The number of requests to Home route",
	})

	PathLength = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Name:      "request_path_lengths",
		Help:      "Groups the lengths of keys in buckets",
		Buckets:   []float64{0, 5, 10, 15, 20, 40, 60, 80, 100, 200, 400, 600, 800, 1000},
	}, []string{"LabelStatus"})

	prometheus.MustRegister(HomeHandlerCounter)
	prometheus.MustRegister(PathLength)

}
