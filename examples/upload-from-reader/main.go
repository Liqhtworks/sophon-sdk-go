// Stream from an io.Reader source (no os.File required) using the
// streaming uploads adapter — chunks go to the wire as bytes.NewReader,
// not via os.CreateTemp.
//
//	go run ./examples/upload-from-reader/ ./some.mp4
//
// The source is read once, in order. UploadFile requires io.ReaderAt so
// it can ReadAt across part boundaries from concurrent workers — if your
// real source is a streaming HTTP body, buffer to disk or to a
// bytes.Reader first.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sophon "github.com/Liqhtworks/sophon-sdk-go"
	"github.com/Liqhtworks/sophon-sdk-go/helpers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: upload-from-reader <path>")
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

	// Pull the file fully into memory and wrap it as a *bytes.Reader.
	// This is the "small file" buffered path; for streaming a multi-GB
	// source see the Daisy integration guide section 3.
	f, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	buf, err := io.ReadAll(f)
	_ = f.Close()
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(buf)

	mimeType := "video/mp4"
	if strings.EqualFold(filepath.Ext(srcPath), ".mov") {
		mimeType = "video/quicktime"
	}

	uploads := helpers.NewStreamingUploadsClient(client)
	res, err := helpers.UploadFile(
		context.Background(), uploads, reader, int64(len(buf)),
		filepath.Base(srcPath), mimeType,
		helpers.UploadFileOptions{
			Concurrency: 4,
			PartTimeout: 0, // use default 60s
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("upload complete: id=%s bytes=%d sha256=%s\n",
		res.UploadID, res.Bytes, res.Sha256)
}
