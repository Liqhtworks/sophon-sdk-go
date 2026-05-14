# sophon-sdk-go

> **Alpha:** This SDK is in alpha. Please report any bugs or errors by opening an [issue](https://github.com/Liqhtworks/sophon-sdk-go/issues).

Official Go SDK for the [SOPHON Encoding API](https://sophon.rs).

> **This package is generated.** Source lives in [Liqhtworks/sophon-api](https://github.com/Liqhtworks/sophon-api) (`api/openapi.yaml` + `api/sdk/helpers/go/`). Do not edit files in this repository by hand — changes are overwritten on every release.

## Install

```bash
go get github.com/Liqhtworks/sophon-sdk-go@latest
```

Requires Go 1.23+.

## Get an API key

1. Sign in at <https://sophon.rs/account/general>.
2. In **API keys**, create a key for your server-side integration.
3. Copy the `xt_live_...` token when it is shown. It is only shown once.
4. Store it as an environment variable:

```bash
export SOPHON_API_KEY=xt_live_...
export SOPHON_BASE_URL=https://api.liqhtworks.xyz
```

Keep API keys on the server. Do not ship them in client apps, public repos,
logs, or analytics events.

## Quick start

The SDK ships generated transport types and a `helpers/` subpackage with
ergonomic wrappers (chunked upload, job polling, webhook signature
verification). The two `New*Client` functions bridge the generated client
to the helpers' interfaces in one line each.

This is the smallest complete server-side flow: upload a local video, create an
encode job, wait for completion, and download the MP4 output.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    sophon "github.com/Liqhtworks/sophon-sdk-go"
    "github.com/Liqhtworks/sophon-sdk-go/helpers"
)

func main() {
    inputPath := "./source.mov"
    if len(os.Args) > 1 {
        inputPath = os.Args[1]
    }

    apiKey := os.Getenv("SOPHON_API_KEY")
    if apiKey == "" { panic("SOPHON_API_KEY is required") }

    baseURL := os.Getenv("SOPHON_BASE_URL")
    if baseURL == "" {
        baseURL = "https://api.liqhtworks.xyz"
    }

    cfg := sophon.NewConfiguration()
    cfg.Servers = sophon.ServerConfigurations{{URL: baseURL}}
    cfg.AddDefaultHeader("Authorization", "Bearer "+apiKey)
    client := sophon.NewAPIClient(cfg)

    ctx := context.Background()

    uploads   := helpers.NewUploadsClient(client.UploadsAPI)
    jobs      := helpers.NewJobsClient(client.JobsAPI)
    downloads := helpers.NewDownloadsClient(client)

    // 1. Upload a file (chunked, concurrent, resumable).
    reader, size, closer, err := helpers.OpenFileForUpload(inputPath)
    if err != nil { panic(err) }
    defer closer()

    mimeType := "video/mp4"
    if strings.EqualFold(filepath.Ext(inputPath), ".mov") {
        mimeType = "video/quicktime"
    }
    upload, err := helpers.UploadFile(
        ctx, uploads, reader, size, filepath.Base(inputPath), mimeType,
        helpers.UploadFileOptions{
            Concurrency: 4,
            OnProgress: func(p helpers.UploadProgress) {
                fmt.Printf("%d/%d parts\n", p.PartsDone, p.PartsTotal)
            },
        },
    )
    if err != nil { panic(err) }

    // 2. Start an encode.
    job, err := helpers.CreateJob(ctx, jobs,
        helpers.JobSource.Upload(upload.UploadID),
        sophon.SOPHON_ESPRESSO,
        helpers.CreateJobOptions{})
    if err != nil { panic(err) }

    // 3. Wait for it to finish.
    final, err := helpers.WaitForJob(ctx, jobs, job.ID, helpers.WaitForJobOptions{
        Timeout: 30 * time.Minute,
        OnProgress: func(j *helpers.Job) {
            fmt.Printf("job %s: %s\n", j.ID, j.Status)
        },
    })
    if err != nil { panic(err) }
    if final.Status != sophon.COMPLETED {
        panic("job ended in " + string(final.Status))
    }

    // 4. Download the encoded output.
    n, err := helpers.DownloadOutputToFile(ctx, downloads, final.ID, "sophon-output.mp4")
    if err != nil { panic(err) }
    fmt.Printf("wrote sophon-output.mp4 (%d bytes)\n", n)
}
```

For a runnable copy of this flow, see
[`examples/encode-file`](./examples/encode-file).

### Profile choice

Use `sophon-auto` for production unless you need deterministic encoder
settings. The quickstart uses `sophon-espresso` because it is the fastest
smoke-test profile and always produces a new encoded output.

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
| `helpers.NewUploadsClient(client.UploadsAPI)` | Bridge generated `*UploadsAPIService` to `UploadsClient`. Stages each chunk through `os.CreateTemp` (constrained by the OpenAPI generator's `*os.File` body type). |
| `helpers.NewStreamingUploadsClient(client)` | Same surface as `NewUploadsClient` but streams part bodies directly from memory via `bytes.NewReader` — no tempfile per chunk. Preferred for large uploads / Windows. |
| `helpers.NewJobsClient(client.JobsAPI)` | Bridge generated `*JobsAPIService` to `JobsClient`. |
| `helpers.NewDownloadsClient(client)` | Bridge the SDK client to `DownloadsClient` for output downloads. |
| `helpers.UploadFile(ctx, uploads, reader, size, name, mime, opts)` | Slice into chunks, upload with bounded concurrency, retry transient errors with jittered backoff, resume from a prior `UploadID`, report progress. `UploadFileOptions.PartTimeout` bounds each part attempt (default 60s); `Retries=0` is now honored literally (pass a negative value for the default 3). |
| `helpers.OpenFileForUpload(path)` | Open a file path and return `(io.ReaderAt, size, closer, err)` ready for `UploadFile`. |
| `helpers.CreateJob(ctx, jobs, source, profile, opts)` | One-call CreateJob wrapper; auto-generates an idempotency key and normalizes nil metadata to `{}`. |
| `helpers.JobSource.Upload(uploadID)` | Typed constructor for `CreateJobRequest.Source` — wraps the oneOf so callers do not type the discriminator. |
| `helpers.WaitForJob(ctx, jobs, id, opts)` | Poll until terminal status (or a caller-specified set). Returns `*JobTerminalError` on `failed`/`canceled` (default set), `*JobTimeoutError` if the deadline elapses. Short `Timeout` values are now bounded by the timeout deadline, not the poll interval. |
| `helpers.DownloadOutput(ctx, downloads, jobID, w)` | Follow the `/v1/jobs/{id}/output` 302, GET the presigned URL, stream into `w`. Returns bytes written. |
| `helpers.DownloadOutputToFile(ctx, downloads, jobID, path)` | Convenience wrapper that creates `path` and streams the output into it. |
| `helpers.VerifyWebhookSignature(body, sig, ts, secret, opts)` | Constant-time HMAC verification + default 5-minute replay window (raw-bytes primitive). |
| `helpers.VerifyWebhookRequest(r, secret, opts)` | One-call `net/http` wrapper: reads the body, pulls the signature headers, verifies, decodes into `*sophon.WebhookDeliveryPayload`. |

### Typed errors

All helper-layer calls (`UploadFile`, `CreateJob`, `WaitForJob`,
`DownloadOutput`, …) return typed errors so callers can branch with
`errors.As` instead of parsing `*sophon.GenericOpenAPIError.Body()`:

| Type | Surfaces | Notes |
|---|---|---|
| `*helpers.AuthenticationError` | HTTP 401 | bad / missing API key |
| `*helpers.PermissionError` | HTTP 403 | key lacks permission |
| `*helpers.NotFoundError` | HTTP 404 | unknown upload / job id |
| `*helpers.ConflictError` | HTTP 409 | idempotency replay with a different body |
| `*helpers.RateLimitError` | HTTP 429 | `.RetryAfter` parsed from `Retry-After` |
| `*helpers.ServerError` | HTTP 5xx | retryable in `UploadFile`'s retry loop |
| `*helpers.NetworkError` | transport failure | retryable (treated as status=0) |

Each type embeds `*helpers.APIError` which exposes `.HTTPStatus()`,
`.Body()` (raw response bytes), and `.Unwrap()` to the underlying
`*sophon.GenericOpenAPIError`.

## Examples

- [`examples/encode-file`](./examples/encode-file) - upload a local video,
  create an encode job, wait for completion, and download the MP4 output.
- [`examples/upload-from-reader`](./examples/upload-from-reader) - stream
  a buffered source through `NewStreamingUploadsClient` (no tempfile per
  chunk).
- [`examples/resume-upload`](./examples/resume-upload) - finish an upload
  that crashed mid-flight by replaying with `UploadFileOptions.UploadID`.
- [`examples/webhook-server`](./examples/webhook-server) - `net/http`
  endpoint that verifies raw request bytes before JSON parsing.

## Runtime support

- Go 1.23+

## Versioning

`github.com/Liqhtworks/sophon-sdk-go` follows
[SemVer](https://semver.org/), with one pre-1.0 caveat: while we are at
`v0.x`, **minor bumps may include breaking changes**. Pin to the v0.1
line until 1.0:

```bash
go get github.com/Liqhtworks/sophon-sdk-go@v0.1
```

Patch releases (`v0.1.x`) are always backward-compatible — they ship
bug fixes, helper-layer improvements, and additive types. Once we cut
`v1.0.0`, regular SemVer applies and breaking changes only land on
major bumps. See [`CHANGELOG.md`](./CHANGELOG.md) for the per-release
log.

## License

Proprietary — see [`LICENSE`](./LICENSE).
