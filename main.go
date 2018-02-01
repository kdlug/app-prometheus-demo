package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "untagged"

// count requests separatelly for each status code
// key/value combinations
var (
	histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "request_duration_seconds",
		Help: "Histogram of request duriation of /hello endpoint",
	},
		[]string{"status"},
	)

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests grouped by status",
		},
		[]string{"status"}, // label
	)
)

// init registers prometheus metrics
func init() {
	prometheus.MustRegister(counter)
	prometheus.MustRegister(histogram)
}

func main() {
	// print Version
	fmt.Println("Version:", version)

	// serve http using promhttp hanlder
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// we use defer function
	// function arguments are initialized imediatelly
	// but all the body is defered

	// we have to initialize status first because it's used in defer function
	var status int

	defer func(start time.Time) {
		// recording time duration with status code
		// first record a time before the requests processing started,
		// then do the actual work, and finally get the current time again,
		// calculate the difference and send that value to the metrics server.
		histogram.With(prometheus.Labels{
			"status": fmt.Sprint(status),
		}).Observe(time.Since(start).Seconds())

		// requests_total{status="200"} 2385
		counter.With(prometheus.Labels{
			"status": fmt.Sprint(status),
		}).Inc()
	}(time.Now()) // we call a function and provide a parameter timeNow()

	status = randomStatus()

	w.WriteHeader(status) // set header
	w.Write(([]byte("Hello, world")))
}

/* Returns a random status code */
func randomStatus() int {
	// random delay
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	status := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusGatewayTimeout,
		http.StatusInternalServerError,
	}

	return status[rand.Intn(len(status))]

}
