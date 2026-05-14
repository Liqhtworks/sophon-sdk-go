package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	sophon "github.com/Liqhtworks/sophon-sdk-go"
)

// Webhook delivery header names. Exported so tests / proxies can rewrite
// them if upstream renames an envelope header without bumping major.
const (
	WebhookSignatureHeader = "X-Turbo-Signature-256"
	WebhookTimestampHeader = "X-Turbo-Timestamp"
)

// WebhookSignatureReason classifies a verification failure.
type WebhookSignatureReason string

const (
	ReasonMissingSignature      WebhookSignatureReason = "missing_signature"
	ReasonMissingTimestamp      WebhookSignatureReason = "missing_timestamp"
	ReasonInvalidTimestamp      WebhookSignatureReason = "invalid_timestamp"
	ReasonReplayWindowExceeded  WebhookSignatureReason = "replay_window_exceeded"
	ReasonBadPrefix             WebhookSignatureReason = "bad_prefix"
	ReasonBadSignatureEncoding  WebhookSignatureReason = "bad_signature_encoding"
	ReasonSignatureMismatch     WebhookSignatureReason = "signature_mismatch"
)

// WebhookSignatureError is returned when a delivery cannot be authenticated.
type WebhookSignatureError struct {
	Reason WebhookSignatureReason
}

func (e *WebhookSignatureError) Error() string {
	return fmt.Sprintf("webhook signature: %s", e.Reason)
}

// VerifyWebhookSignatureOptions tunes the replay window and supplies a
// deterministic clock for tests. Zero values mean "use the default".
type VerifyWebhookSignatureOptions struct {
	// Max acceptable clock skew. The zero value uses the default 5 minutes.
	// Pass any negative duration (e.g. -1) to disable replay enforcement.
	//
	// (Note: the TypeScript/Python helpers treat 0 as "disable" because
	// neither language has a distinct "not set" marker for a numeric field.
	// Go's struct-literal idiom makes the zero value mean "use default"
	// without requiring an explicit pointer, so we diverge on purpose here.)
	ReplayWindow time.Duration
	// Override "now" — used by tests.
	Now func() time.Time
}

// VerifyWebhookSignature returns nil if the delivery is authentic, or a
// *WebhookSignatureError otherwise. `rawBody` must be the raw request bytes
// before any JSON parsing.
//
// SOPHON signs each delivery with HMAC-SHA256 over "{timestamp}.{raw_body}"
// using the per-webhook secret. The hex digest is sent as
// X-Turbo-Signature-256: sha256=<hex>.
func VerifyWebhookSignature(
	rawBody []byte,
	signatureHeader string,
	timestampHeader string,
	secret string,
	opts VerifyWebhookSignatureOptions,
) error {
	if signatureHeader == "" {
		return &WebhookSignatureError{Reason: ReasonMissingSignature}
	}
	if timestampHeader == "" {
		return &WebhookSignatureError{Reason: ReasonMissingTimestamp}
	}

	deliveredTs, err := time.Parse(time.RFC3339Nano, timestampHeader)
	if err != nil {
		deliveredTs, err = time.Parse(time.RFC3339, timestampHeader)
		if err != nil {
			return &WebhookSignatureError{Reason: ReasonInvalidTimestamp}
		}
	}

	window := opts.ReplayWindow
	if window == 0 {
		window = 5 * time.Minute
	}
	if window > 0 {
		nowFn := opts.Now
		if nowFn == nil {
			nowFn = time.Now
		}
		drift := nowFn().Sub(deliveredTs)
		if drift < 0 {
			drift = -drift
		}
		if drift > window {
			return &WebhookSignatureError{Reason: ReasonReplayWindowExceeded}
		}
	}

	if !strings.HasPrefix(signatureHeader, "sha256=") {
		return &WebhookSignatureError{Reason: ReasonBadPrefix}
	}
	deliveredHex := strings.TrimSpace(strings.TrimPrefix(signatureHeader, "sha256="))
	delivered, err := hex.DecodeString(deliveredHex)
	if err != nil {
		return &WebhookSignatureError{Reason: ReasonBadSignatureEncoding}
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestampHeader))
	mac.Write([]byte{'.'})
	mac.Write(rawBody)
	expected := mac.Sum(nil)

	if !hmac.Equal(delivered, expected) {
		return &WebhookSignatureError{Reason: ReasonSignatureMismatch}
	}
	return nil
}

// VerifyWebhookRequest is the one-call wrapper around an inbound *http.Request:
// read the body, pull the SOPHON signature + timestamp headers, verify the
// HMAC, and unmarshal into a WebhookDeliveryPayload. The request body is
// replaced with an io.NopCloser over the bytes that were verified so
// downstream handlers can still ReadAll if they want.
//
//	func handle(w http.ResponseWriter, r *http.Request) {
//	    payload, err := helpers.VerifyWebhookRequest(r, secret, helpers.VerifyWebhookSignatureOptions{})
//	    if err != nil { http.Error(w, "unauthorized", 401); return }
//	    // … use payload …
//	}
func VerifyWebhookRequest(
	r *http.Request,
	secret string,
	opts VerifyWebhookSignatureOptions,
) (*sophon.WebhookDeliveryPayload, error) {
	if r == nil || r.Body == nil {
		return nil, fmt.Errorf("sophon: nil http request")
	}
	body, err := io.ReadAll(r.Body)
	_ = r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("sophon: read webhook body: %w", err)
	}
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	if err := VerifyWebhookSignature(
		body,
		r.Header.Get(WebhookSignatureHeader),
		r.Header.Get(WebhookTimestampHeader),
		secret,
		opts,
	); err != nil {
		return nil, err
	}

	var payload sophon.WebhookDeliveryPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("sophon: decode webhook payload: %w", err)
	}
	return &payload, nil
}
