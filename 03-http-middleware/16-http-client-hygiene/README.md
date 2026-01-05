# Kata 16: The HTTP Client Hygiene Wrapper
**Target Idioms:** `net/http` Transport Reuse, Timeouts, Context-First APIs, Response Body Draining  
**Difficulty:** 🔴 Advanced

## 🧠 The "Why"
“Works locally” HTTP code in Go often fails in prod because people:
- use `http.DefaultClient` with no timeouts,
- create a new client/transport per request (connection churn),
- forget to close bodies (leaks + no keep-alive reuse),
- don’t drain bodies (prevents connection reuse).

This kata is about building a small internal SDK the **Go way**.

## 🎯 The Scenario
Your service calls a downstream API that sometimes returns large error bodies and sometimes hangs.
You need:
- strict timeouts,
- proper cancellation,
- safe connection reuse,
- structured logs.

## 🛠 The Challenge
Implement:
- `type APIClient struct { ... }`
- `func (c *APIClient) GetJSON(ctx context.Context, url string, out any) error`

### 1. Functional Requirements
- [ ] Use `http.NewRequestWithContext`.
- [ ] Decode JSON on 2xx responses into `out`.
- [ ] On non-2xx: read up to N bytes of body and return an error including status code.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must NOT** use `http.DefaultClient`.
- [ ] **Must** configure timeouts (`Client.Timeout` and/or transport-level timeouts).
- [ ] **Must** reuse a single `Transport` (connection pooling).
- [ ] **Must** `defer resp.Body.Close()`.
- [ ] **Must** drain (at least partially) error bodies to allow connection reuse.
- [ ] Use `slog` with fields: method, url, status, latency.

## 🧪 Self-Correction (Test Yourself)
- **If connections spike under load:** you probably rebuild transports.
- **If keep-alives don’t work:** you likely didn’t drain/close body.
- **If hangs occur:** you likely lack correct timeout configuration.

## 📚 Resources
- https://go.dev/src/net/http/client.go
- https://go.dev/src/net/http/transport.go
- https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
