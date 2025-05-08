package main

import (
	"fmt"
	"net/http"
)

func main() {
	initMetrics()
	initAlerting()

	http.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`Welcome! This demo application showcases the capabilities of VictoriaMetrics.
<br />
<br />
Metrics exposed at <a href="/metrics">/metrics</a> endpoint.
<br />
<br />
The /alerting/webhook endpoint is used to receive webhooks from the alert manager.
The <a href="/alerting/receivedWebhooks">/alerting/receivedWebhooks</a> endpoint is used to list the received webhooks.
`))
	}))

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
