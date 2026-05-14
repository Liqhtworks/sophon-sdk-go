package helpers

// Adapters bridging the generated openapi-generator builder API
// (`*UploadsAPIService`, `*JobsAPIService`, `*DownloadsAPIService`) to the
// small interfaces the helpers consume.
//
// This file is shipped as `adapters.go.tmpl` in the helpers source tree
// (api/sdk/helpers/go/helpers/ in sophon-api) because it imports the
// generated parent package, which only exists once the helpers are
// spliced into the published SDK at github.com/Liqhtworks/sophon-sdk-go.
// publish.sh renames the file to `adapters.go` during the splice step
// so it compiles in the published artifact.
//
// In the standalone helper module (sophon-sdk-helpers-source) the file
// keeps the .tmpl extension and Go ignores it — so unit tests for the
// helpers' internal logic still run without the parent package.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	sophon "github.com/Liqhtworks/sophon-sdk-go"
)

// NewUploadsClient wraps the generated `*sophon.UploadsAPIService` so it
// implements `UploadsClient`. Pass `client.UploadsAPI` directly:
//
//	uploads := helpers.NewUploadsClient(client.UploadsAPI)
//	res, err := helpers.UploadFile(ctx, uploads, reader, size, name, mime, opts)
func NewUploadsClient(svc *sophon.UploadsAPIService) UploadsClient {
	return uploadsClientAdapter{svc: svc}
}

type uploadsClientAdapter struct {
	svc *sophon.UploadsAPIService
}

func (a uploadsClientAdapter) CreateUpload(
	ctx context.Context, fileName, mimeType string, fileSize int64, idempotencyKey string,
) (*UploadSession, error) {
	req := sophon.CreateUploadRequest{
		FileName: fileName,
		FileSize: fileSize,
		MimeType: mimeType,
	}
	resp, httpResp, err := a.svc.CreateUpload(ctx).
		IdempotencyKey(idempotencyKey).
		CreateUploadRequest(req).
		Execute()
	if err != nil {
		return nil, classifyError(httpResp, err)
	}
	return &UploadSession{
		ID:          resp.GetId(),
		ChunkSize:   resp.GetChunkSize(),
		TotalChunks: resp.GetTotalChunks(),
	}, nil
}

func (a uploadsClientAdapter) UploadPart(
	ctx context.Context, id string, partNumber int32, body []byte,
) error {
	// OpenAPI Generator's Go output models `format: binary` request bodies
	// as `*os.File`. Customers passing in-memory chunks (the common case
	// when chunking a Blob/io.ReaderAt) need a real file. Stage the chunk
	// in a unique temp file, pass it, and clean up unconditionally.
	f, err := os.CreateTemp("", fmt.Sprintf("sophon-upload-%s-%d-*.bin", id, partNumber))
	if err != nil {
		return fmt.Errorf("create temp for part %d: %w", partNumber, err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	if _, err := f.Write(body); err != nil {
		return fmt.Errorf("write temp for part %d: %w", partNumber, err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seek temp for part %d: %w", partNumber, err)
	}
	_, httpResp, err := a.svc.UploadPart(ctx, id, partNumber).Body(f).Execute()
	if err != nil {
		return classifyError(httpResp, err)
	}
	return nil
}

func (a uploadsClientAdapter) CompleteUpload(
	ctx context.Context, id, idempotencyKey string,
) (*UploadResult, error) {
	resp, httpResp, err := a.svc.CompleteUpload(ctx, id).
		IdempotencyKey(idempotencyKey).
		Execute()
	if err != nil {
		return nil, classifyError(httpResp, err)
	}
	return &UploadResult{
		UploadID: resp.GetId(),
		Sha256:   resp.GetSha256(),
		Bytes:    resp.GetBytes(),
	}, nil
}

func (a uploadsClientAdapter) GetUpload(
	ctx context.Context, id string,
) (*UploadStatus, error) {
	resp, httpResp, err := a.svc.GetUpload(ctx, id).Execute()
	if err != nil {
		return nil, classifyError(httpResp, err)
	}
	return &UploadStatus{
		ID:             resp.GetId(),
		TotalChunks:    int64(resp.GetTotalChunks()),
		ReceivedChunks: resp.GetReceivedChunks(),
	}, nil
}

// NewStreamingUploadsClient wraps a full `*sophon.APIClient` and returns an
// `UploadsClient` that streams part bodies directly from memory instead of
// staging each chunk through `os.CreateTemp`. Behaves identically to
// `NewUploadsClient` except that UploadPart issues a plain `PUT` with
// `bytes.NewReader(body)` against `/v1/uploads/{id}/parts/{n}` — no tempfile
// per chunk. Use this for large uploads or on Windows where AV scanners
// fight temp dir churn.
//
//	uploads := helpers.NewStreamingUploadsClient(client)
//	res, err := helpers.UploadFile(ctx, uploads, reader, size, name, mime, opts)
func NewStreamingUploadsClient(client *sophon.APIClient) UploadsClient {
	return &streamingUploadsAdapter{
		client: client,
		svc:    client.UploadsAPI,
	}
}

type streamingUploadsAdapter struct {
	client *sophon.APIClient
	svc    *sophon.UploadsAPIService
}

func (a *streamingUploadsAdapter) CreateUpload(
	ctx context.Context, fileName, mimeType string, fileSize int64, idempotencyKey string,
) (*UploadSession, error) {
	return uploadsClientAdapter{svc: a.svc}.CreateUpload(ctx, fileName, mimeType, fileSize, idempotencyKey)
}

func (a *streamingUploadsAdapter) CompleteUpload(
	ctx context.Context, id, idempotencyKey string,
) (*UploadResult, error) {
	return uploadsClientAdapter{svc: a.svc}.CompleteUpload(ctx, id, idempotencyKey)
}

func (a *streamingUploadsAdapter) GetUpload(
	ctx context.Context, id string,
) (*UploadStatus, error) {
	return uploadsClientAdapter{svc: a.svc}.GetUpload(ctx, id)
}

func (a *streamingUploadsAdapter) UploadPart(
	ctx context.Context, id string, partNumber int32, body []byte,
) error {
	cfg := a.client.GetConfig()
	base, err := cfg.ServerURLWithContext(ctx, "UploadsAPIService.UploadPart")
	if err != nil {
		return err
	}
	u := base + "/v1/uploads/" + url.PathEscape(id) +
		"/parts/" + url.PathEscape(fmt.Sprintf("%d", partNumber))

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.ContentLength = int64(len(body))
	for k, v := range cfg.DefaultHeader {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Accept", "application/json")
	if cfg.UserAgent != "" {
		req.Header.Set("User-Agent", cfg.UserAgent)
	}

	hc := cfg.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return classifyError(nil, err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		// GenericOpenAPIError fields are unexported, so we cannot stuff
		// the raw response body into the .Underlying wrapper from this
		// package. classifyError surfaces the status text + headers,
		// which is enough for isRetryable / errors.As branching.
		return classifyError(resp, fmt.Errorf("upload part %d: %s", partNumber, resp.Status))
	}
	return nil
}

// NewJobsClient wraps the generated `*sophon.JobsAPIService` so it
// implements `JobsClient`. Pass `client.JobsAPI` directly:
//
//	jobs := helpers.NewJobsClient(client.JobsAPI)
//	final, err := helpers.WaitForJob(ctx, jobs, jobID, helpers.WaitForJobOptions{})
func NewJobsClient(svc *sophon.JobsAPIService) JobsClient {
	return jobsClientAdapter{svc: svc}
}

type jobsClientAdapter struct {
	svc *sophon.JobsAPIService
}

func (a jobsClientAdapter) GetJob(ctx context.Context, id string) (*Job, error) {
	resp, httpResp, err := a.svc.GetJob(ctx, id).Execute()
	if err != nil {
		return nil, classifyError(httpResp, err)
	}
	job := &Job{
		ID:     resp.GetId(),
		Status: resp.GetStatus(),
	}
	if e := resp.GetError(); e != "" {
		job.Error = e
	}
	return job, nil
}

func (a jobsClientAdapter) CreateJob(
	ctx context.Context, source sophon.UploadJobSource, profile sophon.JobProfile,
	idempotencyKey string, output *sophon.CreateJobOutputOptions, metadata map[string]interface{},
) (*Job, error) {
	req := sophon.CreateJobRequest{
		Source:   source,
		Profile:  profile,
		Output:   output,
		Metadata: metadata,
	}
	resp, httpResp, err := a.svc.CreateJob(ctx).
		IdempotencyKey(idempotencyKey).
		CreateJobRequest(req).
		Execute()
	if err != nil {
		return nil, classifyError(httpResp, err)
	}
	job := &Job{
		ID:     resp.GetId(),
		Status: resp.GetStatus(),
	}
	if e := resp.GetError(); e != "" {
		job.Error = e
	}
	return job, nil
}

// DownloadsClient is the transport surface the download helper needs.
type DownloadsClient interface {
	// GetOutputURL fetches a presigned download URL for a completed job.
	GetOutputURL(ctx context.Context, jobID string) (string, error)
}

// NewDownloadsClient wraps the SDK's APIClient so DownloadOutput can issue
// the GET /v1/jobs/{id}/output redirect probe and the subsequent presigned
// GET using the configuration's HTTPClient and default Authorization header.
func NewDownloadsClient(client *sophon.APIClient) DownloadsClient {
	return &downloadsClientAdapter{client: client}
}

type downloadsClientAdapter struct {
	client *sophon.APIClient
}

func (d *downloadsClientAdapter) GetOutputURL(ctx context.Context, jobID string) (string, error) {
	cfg := d.client.GetConfig()
	base, err := cfg.ServerURLWithContext(ctx, "JobsAPIService.GetJobOutput")
	if err != nil {
		base = "https://api.liqhtworks.xyz"
	}
	u := base + "/v1/jobs/" + url.PathEscape(jobID) + "/output"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	for k, v := range cfg.DefaultHeader {
		req.Header.Set(k, v)
	}
	if cfg.UserAgent != "" {
		req.Header.Set("User-Agent", cfg.UserAgent)
	}

	hc := cfg.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	// Use a clone that does not follow redirects so we can capture Location.
	noFollow := *hc
	noFollow.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := noFollow.Do(req)
	if err != nil {
		return "", classifyError(nil, err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusFound, http.StatusMovedPermanently, http.StatusSeeOther,
		http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		loc := resp.Header.Get("Location")
		if loc == "" {
			return "", fmt.Errorf("sophon: output redirect missing Location header")
		}
		locURL, perr := url.Parse(loc)
		if perr != nil {
			return "", fmt.Errorf("sophon: invalid output redirect: %w", perr)
		}
		if !locURL.IsAbs() {
			baseURL, _ := url.Parse(base)
			if baseURL != nil {
				locURL = baseURL.ResolveReference(locURL)
			}
		}
		return locURL.String(), nil
	}
	return "", classifyError(resp, fmt.Errorf("sophon: output endpoint returned %s", resp.Status))
}
