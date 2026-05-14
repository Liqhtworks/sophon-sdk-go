package helpers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	sophon "github.com/Liqhtworks/sophon-sdk-go"
	"github.com/google/uuid"
)

// Job is the subset of the generated JobResponse the helper needs.
type Job struct {
	ID     string
	Status sophon.JobStatus
	Error  string
}

// JobsClient is the transport surface the job waiter and creator require.
type JobsClient interface {
	GetJob(ctx context.Context, id string) (*Job, error)
	CreateJob(
		ctx context.Context,
		source sophon.UploadJobSource,
		profile sophon.JobProfile,
		idempotencyKey string,
		output *sophon.CreateJobOutputOptions,
		metadata map[string]interface{},
	) (*Job, error)
}

// TerminalStatuses is the set of job statuses that stop polling by default.
var TerminalStatuses = map[sophon.JobStatus]struct{}{
	sophon.COMPLETED: {},
	sophon.FAILED:    {},
	sophon.CANCELED:  {},
}

// JobTerminalError is returned when a job ends in failed/canceled and the
// caller was waiting for the default terminal set.
type JobTerminalError struct {
	Job *Job
}

func (e *JobTerminalError) Error() string {
	if e.Job.Error != "" {
		return e.Job.Error
	}
	return fmt.Sprintf("job %s ended in status %s", e.Job.ID, e.Job.Status)
}

// JobTimeoutError is returned when the poll deadline elapses.
type JobTimeoutError struct {
	JobID    string
	WaitedMs int64
}

func (e *JobTimeoutError) Error() string {
	return fmt.Sprintf("job %s did not finish within %dms", e.JobID, e.WaitedMs)
}

// WaitForJobOptions tunes WaitForJob. Zero values mean "use the default".
type WaitForJobOptions struct {
	// Poll until the job reaches any of these statuses. Defaults to the
	// three terminal statuses (completed, failed, canceled).
	Until []sophon.JobStatus
	// Initial poll interval. Default 1s.
	PollMin time.Duration
	// Cap on poll interval. Default 15s.
	PollMax time.Duration
	// Exponential backoff multiplier per poll. Default 1.5.
	PollBackoff float64
	// Hard timeout. Default 1h.
	Timeout time.Duration
	// Optional callback on every poll.
	OnProgress func(*Job)
}

// WaitForJob polls api.GetJob until the job hits a terminal status (or the
// requested Until set), then returns the final job. Returns a *JobTerminalError
// on failed/canceled when using the default terminal set, or a *JobTimeoutError
// if the deadline elapses.
func WaitForJob(ctx context.Context, api JobsClient, jobID string, opts WaitForJobOptions) (*Job, error) {
	waitSet := TerminalStatuses
	customUntil := false
	if len(opts.Until) > 0 {
		customUntil = true
		waitSet = make(map[sophon.JobStatus]struct{}, len(opts.Until))
		for _, s := range opts.Until {
			waitSet[s] = struct{}{}
		}
	}

	pollMin := opts.PollMin
	if pollMin <= 0 {
		pollMin = time.Second
	}
	pollMax := opts.PollMax
	if pollMax <= 0 {
		pollMax = 15 * time.Second
	}
	backoff := opts.PollBackoff
	if backoff <= 0 {
		backoff = 1.5
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = time.Hour
	}

	start := time.Now()
	deadline := start.Add(timeout)
	interval := pollMin

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		elapsed := time.Since(start)
		if elapsed > timeout {
			return nil, &JobTimeoutError{JobID: jobID, WaitedMs: elapsed.Milliseconds()}
		}

		job, err := api.GetJob(ctx, jobID)
		if err != nil {
			return nil, fmt.Errorf("get job: %w", err)
		}
		if opts.OnProgress != nil {
			opts.OnProgress(job)
		}

		if _, match := waitSet[job.Status]; match {
			if !customUntil && (job.Status == sophon.FAILED || job.Status == sophon.CANCELED) {
				return nil, &JobTerminalError{Job: job}
			}
			return job, nil
		}

		// Bound the sleep on the timeout deadline so a short Timeout does
		// not get overrun by up to one full poll interval.
		sleep := interval
		if remaining := time.Until(deadline); remaining < sleep {
			if remaining <= 0 {
				return nil, &JobTimeoutError{JobID: jobID, WaitedMs: time.Since(start).Milliseconds()}
			}
			sleep = remaining
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(sleep):
		}
		next := time.Duration(math.Ceil(float64(interval) * backoff))
		if next > pollMax {
			next = pollMax
		}
		interval = next
	}
}

// CreateJobOptions tunes CreateJob. Zero values mean "use the default".
type CreateJobOptions struct {
	// Auto-generated when empty.
	IdempotencyKey string
	// Optional output settings (container, etc).
	Output *sophon.CreateJobOutputOptions
	// Optional metadata. Normalized to {} when nil so the request body is
	// always a valid map.
	Metadata map[string]interface{}
}

// CreateJob is the one-call wrapper around the generated CreateJob builder.
// It auto-generates an idempotency key when none is supplied and normalizes
// nil metadata to an empty map, matching the Python helper.
//
//	job, err := helpers.CreateJob(ctx, jobs,
//	    helpers.JobSource.Upload(upload.UploadID),
//	    sophon.SOPHON_AUTO,
//	    helpers.CreateJobOptions{})
func CreateJob(
	ctx context.Context,
	api JobsClient,
	source sophon.UploadJobSource,
	profile sophon.JobProfile,
	opts CreateJobOptions,
) (*Job, error) {
	idem := opts.IdempotencyKey
	if idem == "" {
		idem = "idem-" + uuid.NewString()
	}
	meta := opts.Metadata
	if meta == nil {
		meta = map[string]interface{}{}
	}
	return api.CreateJob(ctx, source, profile, idem, opts.Output, meta)
}

// Compile-time proof that our error types implement `error`.
var (
	_ error = (*JobTerminalError)(nil)
	_ error = (*JobTimeoutError)(nil)
	_       = errors.Is
)
