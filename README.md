# sophon-sdk-go

Official Go SDK for the [SOPHON Encoding API](https://liqhtworks.xyz).

> **This package is generated.** Source lives in [Liqhtworks/sophon-api](https://github.com/Liqhtworks/sophon-api) (`api/openapi.yaml` + `api/sdk/helpers/go/`). Do not edit files in this repository by hand — changes are overwritten on every release.

## Install

```bash
go get github.com/Liqhtworks/sophon-sdk-go@latest
```

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "os"

    sophon "github.com/Liqhtworks/sophon-sdk-go"
    "github.com/Liqhtworks/sophon-sdk-go/helpers"
    "github.com/google/uuid"
)

func main() {
    cfg := sophon.NewConfiguration()
    cfg.Servers = sophon.ServerConfigurations{{URL: "https://api.liqhtworks.xyz"}}
    cfg.AddDefaultHeader("Authorization", "Bearer "+os.Getenv("SOPHON_API_KEY"))
    client := sophon.NewAPIClient(cfg)

    ctx := context.Background()

    // 1. Upload a file (chunked, concurrent, resumable).
    reader, size, closer, err := helpers.OpenFileForUpload("/path/to/source.mov")
    if err != nil { panic(err) }
    defer closer()

    uploadsAdapter := // wrap client.UploadsAPI to satisfy helpers.UploadsClient
                      // (see api/sdk/helpers/go/helpers/uploads.go for the interface)

    res, err := helpers.UploadFile(ctx, uploadsAdapter, reader, size,
        "source.mov", "video/quicktime", helpers.UploadFileOptions{})
    if err != nil { panic(err) }

    // 2. Start an encode.
    idempotencyKey := uuid.NewString()
    job, _, err := client.JobsAPI.CreateJob(ctx).
        IdempotencyKey(idempotencyKey).
        CreateJobRequest(sophon.CreateJobRequest{
            Source:  sophon.CreateJobRequestSource{ /* UploadID: res.UploadID */ },
            Profile: "sophon-auto",
        }).
        Execute()
    if err != nil { panic(err) }

    // 3. Wait for it to finish.
    done, err := helpers.WaitForJob(ctx, jobsAdapter, job.Id, helpers.WaitForJobOptions{})
    if err != nil { panic(err) }
    fmt.Println(done.Status)
}
```

## Webhook verification

```go
import "github.com/Liqhtworks/sophon-sdk-go/helpers"

// In your webhook handler — rawBody MUST be the raw request bytes.
err := helpers.VerifyWebhookSignature(
    rawBody,
    r.Header.Get("X-Turbo-Signature-256"),
    r.Header.Get("X-Turbo-Timestamp"),
    os.Getenv("SOPHON_WEBHOOK_SECRET"),
    helpers.VerifyWebhookSignatureOptions{},
)
```

## Runtime support

- Go 1.22+

## License

Proprietary — see [`LICENSE`](./LICENSE).
