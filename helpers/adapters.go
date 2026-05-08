package helpers

// Adapters bridging the generated openapi-generator builder API
// (`*UploadsAPIService`, `*JobsAPIService`) to the small interfaces the
// helpers consume (`UploadsClient`, `JobsClient`).
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
	"context"
	"fmt"
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
	resp, _, err := a.svc.CreateUpload(ctx).
		IdempotencyKey(idempotencyKey).
		CreateUploadRequest(req).
		Execute()
	if err != nil {
		return nil, err
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
	_, _, err = a.svc.UploadPart(ctx, id, partNumber).Body(f).Execute()
	return err
}

func (a uploadsClientAdapter) CompleteUpload(
	ctx context.Context, id, idempotencyKey string,
) (*UploadResult, error) {
	resp, _, err := a.svc.CompleteUpload(ctx, id).
		IdempotencyKey(idempotencyKey).
		Execute()
	if err != nil {
		return nil, err
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
	resp, _, err := a.svc.GetUpload(ctx, id).Execute()
	if err != nil {
		return nil, err
	}
	return &UploadStatus{
		ID:             resp.GetId(),
		TotalChunks:    int64(resp.GetTotalChunks()),
		ReceivedChunks: resp.GetReceivedChunks(),
	}, nil
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
	resp, _, err := a.svc.GetJob(ctx, id).Execute()
	if err != nil {
		return nil, err
	}
	job := &Job{
		ID:     resp.GetId(),
		Status: string(resp.GetStatus()),
	}
	// `Error` is `*string` on the generated model — dereference if set.
	if e := resp.GetError(); e != "" {
		job.Error = e
	}
	return job, nil
}
