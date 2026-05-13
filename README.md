# sophon-sdk-go

> **Alpha:** This SDK is in alpha. Please report any bugs or errors by opening an [issue](https://github.com/Liqhtworks/sophon-sdk-go/issues).

Official Go SDK for the [SOPHON Encoding API](https://sophon.rs).

> **This package is generated.** Source lives in [Liqhtworks/sophon-api](https://github.com/Liqhtworks/sophon-api) (`api/openapi.yaml` + `api/sdk/helpers/go/`). Do not edit files in this repository by hand — changes are overwritten on every release.

## Install

```bash
go get github.com/Liqhtworks/sophon-sdk-go@latest
```

Requires Go 1.22+.

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
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
    "time"

    sophon "github.com/Liqhtworks/sophon-sdk-go"
    "github.com/Liqhtworks/sophon-sdk-go/helpers"
    "github.com/google/uuid"
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

    uploads := helpers.NewUploadsClient(client.UploadsAPI)
    jobs    := helpers.NewJobsClient(client.JobsAPI)

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
    idempotencyKey := uuid.NewString()
    job, _, err := client.JobsAPI.CreateJob(ctx).
        IdempotencyKey(idempotencyKey).
        CreateJobRequest(sophon.CreateJobRequest{
            Source:  helpers.JobSource.Upload(upload.UploadID),
            Profile: sophon.SOPHON_ESPRESSO,
        }).
        Execute()
    if err != nil { panic(err) }

    // 3. Wait for it to finish.
    final, err := helpers.WaitForJob(ctx, jobs, job.GetId(), helpers.WaitForJobOptions{
        Timeout: 30 * time.Minute,
        OnProgress: func(j *helpers.Job) {
            fmt.Printf("job %s: %s\n", j.ID, j.Status)
        },
    })
    if err != nil { panic(err) }
    if final.Status != "completed" { panic("job ended in " + final.Status) }

    // 4. Download the encoded output.
    if err := downloadOutput(baseURL, apiKey, final.ID, "sophon-output.mp4"); err != nil {
        panic(err)
    }
    fmt.Println("wrote sophon-output.mp4")
}

func downloadOutput(baseURL, apiKey, jobID, outputPath string) error {
    client := &http.Client{
        CheckRedirect: func(*http.Request, []*http.Request) error {
            return http.ErrUseLastResponse
        },
        Timeout: 30 * time.Second,
    }

    req, err := http.NewRequest("GET", baseURL+"/v1/jobs/"+jobID+"/output", nil)
    if err != nil { return err }
    req.Header.Set("Authorization", "Bearer "+apiKey)

    res, err := client.Do(req)
    if err != nil { return err }
    defer res.Body.Close()
    if res.StatusCode != http.StatusFound {
        return fmt.Errorf("expected output redirect, got %d", res.StatusCode)
    }

    location := res.Header.Get("Location")
    if location == "" { return fmt.Errorf("missing output redirect") }

    downloadURL, err := resolveURL(baseURL, location)
    if err != nil { return err }

    dl, err := http.Get(downloadURL)
    if err != nil { return err }
    defer dl.Body.Close()
    if dl.StatusCode < 200 || dl.StatusCode >= 300 {
        return fmt.Errorf("download failed: %d", dl.StatusCode)
    }

    out, err := os.Create(outputPath)
    if err != nil { return err }
    defer out.Close()
    _, err = io.Copy(out, dl.Body)
    return err
}

func resolveURL(baseURL, location string) (string, error) {
    loc, err := url.Parse(location)
    if err != nil { return "", err }
    if loc.IsAbs() { return loc.String(), nil }
    base, err := url.Parse(baseURL)
    if err != nil { return "", err }
    return base.ResolveReference(loc).String(), nil
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
| `helpers.NewUploadsClient(client.UploadsAPI)` | Bridge generated `*UploadsAPIService` to `UploadsClient`. |
| `helpers.NewJobsClient(client.JobsAPI)` | Bridge generated `*JobsAPIService` to `JobsClient`. |
| `helpers.UploadFile(ctx, uploads, reader, size, name, mime, opts)` | Slice into chunks, upload with bounded concurrency, retry transient errors with jittered backoff, resume from a prior `UploadID`, report progress. |
| `helpers.OpenFileForUpload(path)` | Open a file path and return `(io.ReaderAt, size, closer, err)` ready for `UploadFile`. |
| `helpers.WaitForJob(ctx, jobs, id, opts)` | Poll until terminal status (or a caller-specified set). Returns `*JobTerminalError` on `failed`/`canceled` (default set), `*JobTimeoutError` if the deadline elapses. |
| `helpers.VerifyWebhookSignature(body, sig, ts, secret, opts)` | Constant-time HMAC verification + default 5-minute replay window. |

## Examples

- [`examples/encode-file`](./examples/encode-file) - upload a local video,
  create an encode job, wait for completion, and download the MP4 output.
- [`examples/webhook-server`](./examples/webhook-server) - `net/http`
  endpoint that verifies raw request bytes before JSON parsing.

## Runtime support

- Go 1.22+

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
