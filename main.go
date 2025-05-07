package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

var demoCounterTotal = metrics.NewCounter(`demo_counter_total`)

func main() {
	go func() {
		for {
			demoCounterTotal.Inc()
			time.Sleep(randDurationInRange(time.Millisecond*100, time.Second*2))
		}
	}()

	http.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Write([]byte(`Welcome! 

This demo app exposes Prometheus metrics in the standard exposition format. Visit the /metrics endpoint.
`))
	}))

	http.Handle("/metrics", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		metrics.WritePrometheus(rw, false)
	}))

	fmt.Println("Starting server on :9100")
	http.ListenAndServe(":9100", nil)
}

func randDurationInRange(min, max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max-min))) + min
}
