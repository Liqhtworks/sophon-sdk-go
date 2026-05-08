---
name: Bug report
about: Something the SDK does that contradicts the docs or its types.
title: "[bug] "
labels: bug
---

## What happened

<!-- What you tried to do, what the SDK did instead. -->

## Reproducer

```go
// Minimum code that reproduces. Strip secrets.
import (
    sophon "github.com/Liqhtworks/sophon-sdk-go"
    "github.com/Liqhtworks/sophon-sdk-go/helpers"
)
// …
```

## Environment

- `github.com/Liqhtworks/sophon-sdk-go` version: `v0.1.x` (from `go list -m github.com/Liqhtworks/sophon-sdk-go`)
- Go version: `go version`
- OS / arch: `…`

## Expected vs. actual

- Expected: `…`
- Actual: `…` (paste error inside a fenced block)

## Anything else

<!-- Logs, X-Request-Id headers from the response, etc. -->
