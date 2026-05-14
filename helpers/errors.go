package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	sophon "github.com/Liqhtworks/sophon-sdk-go"
)

// APIError is the common shape for typed HTTP errors surfaced by the helper
// layer. It implements HTTPStatusCarrier so helpers.isRetryable can classify
// generated-client failures, and wraps the underlying *sophon.GenericOpenAPIError
// so callers can still reach .Body() bytes when needed.
type APIError struct {
	Status     int
	Message    string
	RetryAfter time.Duration
	Underlying *sophon.GenericOpenAPIError
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Underlying != nil {
		return e.Underlying.Error()
	}
	return fmt.Sprintf("sophon api error: status %d", e.Status)
}

func (e *APIError) HTTPStatus() int { return e.Status }

func (e *APIError) Unwrap() error {
	if e.Underlying == nil {
		return nil
	}
	return e.Underlying
}

// Body returns the raw response body bytes from the generated error, if any.
func (e *APIError) Body() []byte {
	if e.Underlying == nil {
		return nil
	}
	return e.Underlying.Body()
}

// Typed errors callers can match with errors.As.
type (
	AuthenticationError struct{ *APIError }
	PermissionError     struct{ *APIError }
	NotFoundError       struct{ *APIError }
	ConflictError       struct{ *APIError }
	RateLimitError      struct{ *APIError }
	ServerError         struct{ *APIError }
	NetworkError        struct{ *APIError }
)

// classifyError inspects the *http.Response returned alongside a generated
// error and produces a typed APIError variant. If resp is nil it returns the
// original error so callers can still distinguish transport errors. err is
// expected to be non-nil.
func classifyError(resp *http.Response, err error) error {
	if err == nil {
		return nil
	}
	var ge *sophon.GenericOpenAPIError
	_ = errors.As(err, &ge)

	if resp == nil {
		// Pure transport failure (DNS, dial, TLS, EOF). Wrap as NetworkError
		// so isRetryable can treat it as retryable.
		return &NetworkError{&APIError{Status: 0, Message: err.Error(), Underlying: ge}}
	}

	base := &APIError{Status: resp.StatusCode, Message: err.Error(), Underlying: ge}
	switch resp.StatusCode {
	case 401:
		return &AuthenticationError{base}
	case 403:
		return &PermissionError{base}
	case 404:
		return &NotFoundError{base}
	case 409:
		return &ConflictError{base}
	case 429:
		base.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"))
		return &RateLimitError{base}
	}
	if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return &ServerError{base}
	}
	if resp.StatusCode >= 400 {
		return base
	}
	return err
}

func parseRetryAfter(h string) time.Duration {
	h = strings.TrimSpace(h)
	if h == "" {
		return 0
	}
	if secs, err := strconv.Atoi(h); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(h); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}
	return 0
}
