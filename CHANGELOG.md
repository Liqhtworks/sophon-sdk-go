# Changelog

All notable changes to `github.com/Liqhtworks/sophon-sdk-go` are
recorded here. The module follows [SemVer](https://semver.org/) — see
`README.md` for the versioning policy applied during the v0.x pre-1.0
phase.

## [0.1.4] — 2026-05-08

- `helpers.UploadJobSource(uploadID)` and the matching constructor on
  `CreateJobRequest.Source` — typed alternative to constructing the
  `oneOf` discriminated union by hand.
- Generated and helper exports tightened so the customer-facing surface
  is reachable from the two documented import paths
  (`github.com/Liqhtworks/sophon-sdk-go` and `.../helpers`).

## [0.1.2] — 2026-04-23

- Per-route idempotency keys in `helpers.UploadFile`. Earlier releases
  reused one key for both `CreateUpload` and `CompleteUpload`; SOPHON
  scopes idempotency keys per route and rejected the second call with
  HTTP 409. Now derives `idem+"/create"` and `idem+"/complete"` from
  the caller's seed so retries still reach the server's idempotent path.

## [0.1.1] — 2026-04-23

- `helpers.NewUploadsClient(*UploadsAPIService)` and
  `helpers.NewJobsClient(*JobsAPIService)` adapters bridge the
  generated builder API to the helpers' small interfaces. Customers
  used to hand-write ~80 lines of glue per program; now it's one line
  each.

## [0.1.0] — 2026-04-23

Initial public release.

- Generated transport (`*sophon.APIClient`, `*JobsAPIService`,
  `*UploadsAPIService`, `*WebhooksAPIService`, `*DownloadsAPIService`,
  `*HealthAPIService`) from the SOPHON OpenAPI spec.
- Hand-written helpers under `github.com/Liqhtworks/sophon-sdk-go/helpers`:
  - `UploadFile` — chunked, concurrent, resumable upload with progress
    callback, bounded retry, and context cancellation.
  - `WaitForJob` — typed terminal-state polling with backoff, timeout,
    and `JobTerminalError` / `JobTimeoutError`.
  - `VerifyWebhookSignature` — constant-time HMAC-SHA256 verification
    with a default replay window.
- `OpenFileForUpload` convenience opener for `(io.ReaderAt, size, closer)`.

[0.1.4]: https://github.com/Liqhtworks/sophon-sdk-go/releases/tag/v0.1.4
[0.1.2]: https://github.com/Liqhtworks/sophon-sdk-go/releases/tag/v0.1.2
[0.1.1]: https://github.com/Liqhtworks/sophon-sdk-go/releases/tag/v0.1.1
[0.1.0]: https://github.com/Liqhtworks/sophon-sdk-go/releases/tag/v0.1.0
