# sophon-sdk-go

Official Go SDK for the [SOPHON Encoding API](https://liqhtworks.xyz).

> **This package is generated.** Source lives in [Liqhtworks/sophon-api](https://github.com/Liqhtworks/sophon-api) (`api/openapi.yaml` + `api/sdk/helpers/go/`). Do not edit files in this repository by hand — changes are overwritten on every release.

## Install

```bash
go get github.com/Liqhtworks/sophon-sdk-go@latest
```

Requires Go 1.22+.

## Quick start

The SDK ships generated transport types and a `helpers/` subpackage with
ergonomic wrappers (chunked upload, job polling, webhook signature
verification). The two `New*Client` functions bridge the generated client
to the helpers' interfaces in one line each.

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

    uploads := helpers.NewUploadsClient(client.UploadsAPI)
    jobs    := helpers.NewJobsClient(client.JobsAPI)

    // 1. Upload a file (chunked, concurrent, resumable).
    reader, size, closer, err := helpers.OpenFileForUpload("/path/to/source.mov")
    if err != nil { panic(err) }
    defer closer()

    upload, err := helpers.UploadFile(
        ctx, uploads, reader, size, "source.mov", "video/quicktime",
        helpers.UploadFileOptions{
            Concurrency: 4,
            OnProgress: func(p helpers.UploadProgress) {
                fmt.Printf("%d/%d parts\n", p.PartsDone, p.PartsTotal)
            },
        },
    )
    if err != nil { panic(err) }

    // 2. Start an encode.
    idempotencyKey := uuid.NewString()
    job, _, err := client.JobsAPI.CreateJob(ctx).
        IdempotencyKey(idempotencyKey).
        CreateJobRequest(sophon.CreateJobRequest{
            Source: sophon.UploadJobSourceAsCreateJobRequestSource(
                sophon.NewUploadJobSource("upload", upload.UploadID),
            ),
            Profile: "sophon-auto",
        }).
        Execute()
    if err != nil { panic(err) }

    // 3. Wait for it to finish.
    final, err := helpers.WaitForJob(ctx, jobs, job.GetId(), helpers.WaitForJobOptions{})
    if err != nil { panic(err) }
    fmt.Println("status:", final.Status)
}
```

## Webhook verification

`helpers.VerifyWebhookSignature` does a constant-time HMAC-SHA256 check
of `"{timestamp}.{raw_body}"` against the `sha256=` header, plus a default
5-minute replay window.

For a complete `net/http` server, see
[`examples/webhook-server`](./examples/webhook-server).

```go
import (
    "io"
    "net/http"
    "os"

    "github.com/Liqhtworks/sophon-sdk-go/helpers"
)

func handle(w http.ResponseWriter, r *http.Request) {
    rawBody, _ := io.ReadAll(r.Body)

    err := helpers.VerifyWebhookSignature(
        rawBody,
        r.Header.Get("X-Turbo-Signature-256"),
        r.Header.Get("X-Turbo-Timestamp"),
        os.Getenv("SOPHON_WEBHOOK_SECRET"),
        helpers.VerifyWebhookSignatureOptions{},
    )
    if err != nil {
        // *helpers.WebhookSignatureError carries .Reason for granular handling.
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // … parse rawBody as a WebhookDeliveryPayload and process …
    w.WriteHeader(http.StatusOK)
}
```

## Helpers reference

| Function | What it does |
|---|---|
| `helpers.NewUploadsClient(client.UploadsAPI)` | Bridge generated `*UploadsAPIService` to `UploadsClient`. |
| `helpers.NewJobsClient(client.JobsAPI)` | Bridge generated `*JobsAPIService` to `JobsClient`. |
| `helpers.UploadFile(ctx, uploads, reader, size, name, mime, opts)` | Slice into chunks, upload with bounded concurrency, retry transient errors with jittered backoff, resume from a prior `UploadID`, report progress. |
| `helpers.OpenFileForUpload(path)` | Open a file path and return `(io.ReaderAt, size, closer, err)` ready for `UploadFile`. |
| `helpers.WaitForJob(ctx, jobs, id, opts)` | Poll until terminal status (or a caller-specified set). Returns `*JobTerminalError` on `failed`/`canceled` (default set), `*JobTimeoutError` if the deadline elapses. |
| `helpers.VerifyWebhookSignature(body, sig, ts, secret, opts)` | Constant-time HMAC verification + default 5-minute replay window. |

## Examples

- [`examples/webhook-server`](./examples/webhook-server) - `net/http`
  endpoint that verifies raw request bytes before JSON parsing.

## Runtime support

- Go 1.22+

## License

Proprietary — see [`LICENSE`](./LICENSE).
