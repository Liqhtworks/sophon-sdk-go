# net/http Webhook Server

```bash
SOPHON_WEBHOOK_SECRET=whsec_... go run ./examples/webhook-server
```

Register `POST /webhooks/sophon` as the webhook endpoint. The handler reads the
raw body first, verifies `X-Turbo-Signature-256`, and only then parses JSON.
