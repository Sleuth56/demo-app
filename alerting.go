package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"
)

var alertingWebhooksMutex sync.Mutex
var alertingWebhooks [][]byte

func initAlerting() {
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
}
