package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/VictoriaMetrics/metrics"
)

var demoAlertFire = atomic.Int64{}

var alertingWebhooksMutex sync.Mutex
var alertingWebhooks [][]byte

func initAlerting() {
	metrics.NewGauge(`demo_alert_firing`, func() float64 {
		return float64(demoAlertFire.Load())
	})

	http.Handle("/alerting/webhook", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %s", err)
			http.Error(rw, "cannot read body", http.StatusBadRequest)
			return
		}

		alertingWebhooksMutex.Lock()
		if len(alertingWebhooks) > 100 {
			alertingWebhooks = alertingWebhooks[1:]
		}
		alertingWebhooks = append(alertingWebhooks, bytes.TrimSpace(b))
		alertingWebhooksMutex.Unlock()

		rw.WriteHeader(http.StatusNoContent)
	}))

	http.Handle("/alerting/receivedWebhooks", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)

		alertingWebhooksMutex.Lock()
		for _, webhook := range alertingWebhooks {
			if _, err := rw.Write(webhook); err != nil {
				log.Printf("Error writing webhook: %s", err)
				break
			}
			if _, err := rw.Write([]byte("\n")); err != nil {
				log.Printf("Error writing webhook: %s", err)
				break
			}
		}
		alertingWebhooksMutex.Unlock()
	}))

	http.Handle("/alerting/fireDemoAlert", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		demoAlertFire.Store(1)
	}))

	http.Handle("/alerting/resolveDemoAlert", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		demoAlertFire.Store(0)
	}))
}
