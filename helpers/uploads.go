// Package helpers wraps the generated SOPHON SDK with higher-level
// ergonomics: chunked upload orchestration, job polling, and webhook
// signature verification.
//
// The helpers operate against small interfaces so they can be tested with
// fakes in isolation. Consumers bridge the generated `*APIService` types to
// these interfaces via the `New*Client` constructors in adapters.go.
package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

// UploadSession is the subset of CreateUploadResponse the helper needs.
type UploadSession struct {
	ID          string
	ChunkSize   int64
	TotalChunks int64
}

// UploadStatus is the subset of UploadStatusResponse the helper needs.
type UploadStatus struct {
	ID             string
	TotalChunks    int64
	ReceivedChunks []int32
}

// UploadResult is the subset of CompleteUploadResponse the helper returns.
type UploadResult struct {
	UploadID string
	Sha256   string
	Bytes    int64
}

// UploadsClient is the transport surface the uploader requires. Wire
// adapters.NewUploadsClient for the generated SDK; mock directly for tests.
type UploadsClient interface {
	CreateUpload(ctx context.Context, fileName, mimeType string, fileSize int64, idempotencyKey string) (*UploadSession, error)
	UploadPart(ctx context.Context, id string, partNumber int32, body []byte) error
	CompleteUpload(ctx context.Context, id, idempotencyKey string) (*UploadResult, error)
	GetUpload(ctx context.Context, id string) (*UploadStatus, error)
}

// UploadProgress is delivered to the OnProgress callback.
type UploadProgress struct {
	BytesUploaded int64
	BytesTotal    int64
	PartsDone     int64
	PartsTotal    int64
}

// UploadFileOptions tunes UploadFile. Zero values mean "use the default".
type UploadFileOptions struct {
	// Resume an existing upload session instead of creating a new one.
	UploadID string
	// Parallel in-flight parts. Default 4.
	Concurrency int
	// Per-part retry count for retryable errors. Default 3.
	Retries int
	// Base retry backoff; doubles each attempt with jitter. Default 500ms.
	RetryBase time.Duration
	// Reused for create+complete. Auto-generated when empty.
	IdempotencyKey string
	// Optional progress callback.
	OnProgress func(UploadProgress)
}

// RetryableError lets callers mark a non-HTTP error as retryable.
type RetryableError struct{ Err error }

func (e *RetryableError) Error() string { return e.Err.Error() }
func (e *RetryableError) Unwrap() error { return e.Err }

// HTTPStatusCarrier is implemented by errors that wrap an HTTP status.
type HTTPStatusCarrier interface{ HTTPStatus() int }

func isRetryable(err error) bool {
	var re *RetryableError
	if errors.As(err, &re) {
		return true
	}
	var sc HTTPStatusCarrier
	if errors.As(err, &sc) {
		s := sc.HTTPStatus()
		return s == 408 || s == 429 || (s >= 500 && s < 600)
	}
	return false
}

// UploadFile slices source into chunks, uploads them with bounded concurrency
// and retries, and finalizes the session. Returns the completed upload's ID,
// sha256, and byte count.
func UploadFile(
	ctx context.Context,
	api UploadsClient,
	source io.ReaderAt,
	sourceSize int64,
	fileName, mimeType string,
	opts UploadFileOptions,
) (*UploadResult, error) {
	concurrency := opts.Concurrency
	if concurrency <= 0 {
		concurrency = 4
	}
	retries := opts.Retries
	if retries < 0 {
		return nil, fmt.Errorf("retries must be >= 0")
	}
	if retries == 0 {
		retries = 3
	}
	retryBase := opts.RetryBase
	if retryBase <= 0 {
		retryBase = 500 * time.Millisecond
	}
	idem := opts.IdempotencyKey
	if idem == "" {
		idem = "idem-" + uuid.NewString()
	}

	var (
		uploadID    string
		chunkSize   int64
		totalChunks int64
		received    = map[int32]struct{}{}
	)
	if opts.UploadID != "" {
		st, err := api.GetUpload(ctx, opts.UploadID)
		if err != nil {
			return nil, fmt.Errorf("get upload: %w", err)
		}
		uploadID = st.ID
		totalChunks = st.TotalChunks
		for _, p := range st.ReceivedChunks {
			received[p] = struct{}{}
		}
		chunkSize = (sourceSize + totalChunks - 1) / totalChunks
	} else {
		s, err := api.CreateUpload(ctx, fileName, mimeType, sourceSize, idem)
		if err != nil {
			return nil, fmt.Errorf("create upload: %w", err)
		}
		uploadID = s.ID
		chunkSize = s.ChunkSize
		totalChunks = s.TotalChunks
	}

	progress := UploadProgress{
		BytesTotal:  sourceSize,
		PartsTotal:  totalChunks,
	}
	for p := range received {
		progress.PartsDone++
		progress.BytesUploaded += partBytes(sourceSize, chunkSize, totalChunks, p)
	}
	var progressMu sync.Mutex
	emitProgress := func() {
		if opts.OnProgress == nil {
			return
		}
		progressMu.Lock()
		snapshot := progress
		progressMu.Unlock()
		opts.OnProgress(snapshot)
	}
	emitProgress()

	var pending []int32
	for i := int32(0); int64(i) < totalChunks; i++ {
		if _, already := received[i]; !already {
			pending = append(pending, i)
		}
	}
	if len(pending) == 0 {
		// Resumed session with nothing left to do; fall through to complete.
	} else {
		workerCount := concurrency
		if workerCount > len(pending) {
			workerCount = len(pending)
		}
		work := make(chan int32)
		errCh := make(chan error, workerCount)
		var wg sync.WaitGroup

		for w := 0; w < workerCount; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for partNum := range work {
					start := int64(partNum) * chunkSize
					end := start + chunkSize
					if end > sourceSize {
						end = sourceSize
					}
					chunk := make([]byte, end-start)
					if _, err := source.ReadAt(chunk, start); err != nil && err != io.EOF {
						errCh <- fmt.Errorf("read part %d: %w", partNum, err)
						return
					}
					if err := withRetry(ctx, retries, retryBase, func() error {
						return api.UploadPart(ctx, uploadID, partNum, chunk)
					}); err != nil {
						errCh <- fmt.Errorf("upload part %d: %w", partNum, err)
						return
					}
					progressMu.Lock()
					progress.PartsDone++
					progress.BytesUploaded += int64(len(chunk))
					progressMu.Unlock()
					emitProgress()
				}
			}()
		}

		go func() {
			defer close(work)
			for _, p := range pending {
				select {
				case <-ctx.Done():
					return
				case work <- p:
				}
			}
		}()

		wg.Wait()
		close(errCh)
		for err := range errCh {
			if err != nil {
				return nil, err
			}
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
	}

	done, err := api.CompleteUpload(ctx, uploadID, idem)
	if err != nil {
		return nil, fmt.Errorf("complete upload: %w", err)
	}
	return done, nil
}

func partBytes(total, chunkSize, totalChunks int64, part int32) int64 {
	if int64(part) < totalChunks-1 {
		return chunkSize
	}
	return total - chunkSize*(totalChunks-1)
}

func withRetry(ctx context.Context, retries int, base time.Duration, fn func() error) error {
	attempt := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := fn()
		if err == nil {
			return nil
		}
		if attempt >= retries || !isRetryable(err) {
			return err
		}
		delay := base<<attempt + time.Duration(rand.Int64N(int64(base)))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		attempt++
	}
}

// OpenFileForUpload is a convenience wrapper for UploadFile over a file path.
func OpenFileForUpload(path string) (io.ReaderAt, int64, func() error, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, nil, err
	}
	return f, info.Size(), f.Close, nil
}
