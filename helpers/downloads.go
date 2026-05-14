package helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadOutput streams an encoded job output to w. It calls
// GET /v1/jobs/{id}/output to obtain a presigned redirect, follows the
// redirect to a public URL, and copies the body. Returns the byte count
// written to w. Errors classify to *NotFoundError when the job is unknown
// and other typed *APIError variants for non-2xx upstream responses.
//
//	downloads := helpers.NewDownloadsClient(client)
//	f, _ := os.Create("out.mp4")
//	defer f.Close()
//	n, err := helpers.DownloadOutput(ctx, downloads, jobID, f)
func DownloadOutput(ctx context.Context, api DownloadsClient, jobID string, w io.Writer) (int64, error) {
	if w == nil {
		return 0, fmt.Errorf("sophon: nil writer")
	}
	signedURL, err := api.GetOutputURL(ctx, jobID)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, signedURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, &NetworkError{&APIError{Status: 0, Message: err.Error()}}
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, classifyError(resp, fmt.Errorf("sophon: download failed: %s", resp.Status))
	}
	return io.Copy(w, resp.Body)
}

// DownloadOutputToFile is a thin wrapper that creates path (overwriting
// any existing file) and streams the output into it. Closes the file on
// return.
func DownloadOutputToFile(ctx context.Context, api DownloadsClient, jobID, path string) (int64, error) {
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return DownloadOutput(ctx, api, jobID, f)
}
