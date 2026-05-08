package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Liqhtworks/sophon-sdk-go/helpers"
)

func main() {
	secret := os.Getenv("SOPHON_WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("SOPHON_WEBHOOK_SECRET is required")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhooks/sophon", func(w http.ResponseWriter, r *http.Request) {
		rawBody, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 2<<20))
		if err != nil {
			http.Error(w, "read body", http.StatusBadRequest)
			return
		}

		err = helpers.VerifyWebhookSignature(
			rawBody,
			r.Header.Get("X-Turbo-Signature-256"),
			r.Header.Get("X-Turbo-Timestamp"),
			secret,
			helpers.VerifyWebhookSignatureOptions{},
		)
		if err != nil {
			if sigErr, ok := err.(*helpers.WebhookSignatureError); ok {
				log.Printf("rejected SOPHON webhook: reason=%s", sigErr.Reason)
			}
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var event map[string]any
		if err := json.Unmarshal(rawBody, &event); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		log.Printf("accepted SOPHON webhook: type=%v id=%v", event["type"], event["id"])
		w.WriteHeader(http.StatusNoContent)
	})

	addr := ":3000"
	fmt.Println("listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
