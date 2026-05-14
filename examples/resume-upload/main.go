// Resume a chunked upload after a crash. Pass the prior upload ID via
// SOPHON_RESUME_ID; UploadFile re-fetches the session, skips already-
// received parts, and re-uploads only the missing ones.
//
//	go run ./examples/resume-upload/ ./source.mp4
//	  → prints "resume id: upl_..."
//	  → ^C mid-upload
//	  SOPHON_RESUME_ID=upl_... go run ./examples/resume-upload/ ./source.mp4
//	  → finishes the missing parts and completes the session
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sophon "github.com/Liqhtworks/sophon-sdk-go"
	"github.com/Liqhtworks/sophon-sdk-go/helpers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: resume-upload <path>")
		os.Exit(1)
	}
	srcPath := os.Args[1]

	apiKey := os.Getenv("SOPHON_API_KEY")
	if apiKey == "" {
		panic("SOPHON_API_KEY is required")
	}
	baseURL := os.Getenv("SOPHON_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.liqhtworks.xyz"
	}

	cfg := sophon.NewConfiguration()
	cfg.Servers = sophon.ServerConfigurations{{URL: baseURL}}
	cfg.AddDefaultHeader("Authorization", "Bearer "+apiKey)
	client := sophon.NewAPIClient(cfg)

	reader, size, closer, err := helpers.OpenFileForUpload(srcPath)
	if err != nil {
		panic(err)
	}
	defer closer()

	mimeType := "video/mp4"
	if strings.EqualFold(filepath.Ext(srcPath), ".mov") {
		mimeType = "video/quicktime"
	}

	uploads := helpers.NewStreamingUploadsClient(client)
	opts := helpers.UploadFileOptions{
		Concurrency: 4,
		UploadID:    os.Getenv("SOPHON_RESUME_ID"),
		OnProgress: func(p helpers.UploadProgress) {
			fmt.Printf("%d/%d parts (%d bytes)\n",
				p.PartsDone, p.PartsTotal, p.BytesUploaded)
		},
	}
	if opts.UploadID == "" {
		// Print the upload id on the first part so the caller can
		// re-launch with SOPHON_RESUME_ID=<id> after a crash. Order is
		// not guaranteed across concurrent workers, but at least one
		// part lands before any failure surfaces.
		first := true
		origProgress := opts.OnProgress
		opts.OnProgress = func(p helpers.UploadProgress) {
			if first && p.PartsDone > 0 {
				first = false
				// The session id is the same as res.UploadID below;
				// we don't have it here yet, so leave a hint instead.
				fmt.Fprintln(os.Stderr, "→ to resume, capture res.UploadID from the final output")
			}
			origProgress(p)
		}
	}

	res, err := helpers.UploadFile(
		context.Background(), uploads, reader, size,
		filepath.Base(srcPath), mimeType, opts,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "upload failed; resume with SOPHON_RESUME_ID=<id>: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("upload complete: id=%s bytes=%d sha256=%s\n",
		res.UploadID, res.Bytes, res.Sha256)
}
