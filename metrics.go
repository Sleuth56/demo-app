package main

import (
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

var demoCounterTotal = metrics.NewCounter(`demo_counter_total`)

func initMetrics() {
	go func() {
		for {
			demoCounterTotal.Inc()
			time.Sleep(time.Second)
		}
	}()

	http.Handle("/metrics", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		metrics.WritePrometheus(rw, false)
	}))
}
