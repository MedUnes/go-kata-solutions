# Kata 07: The Rate-Limited Fan-Out Client
**Target Idioms:** Rate Limiting (`x/time/rate`), Bounded Concurrency (`x/sync/semaphore`), HTTP Client Hygiene, Context Cancellation  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
In many ecosystems, you slap a “rate limit middleware” in front of a thread pool and call it a day. In Go, people often:
- spawn too many goroutines (no backpressure),
- forget per-request cancellation,
- misuse `http.DefaultClient` (timeouts/transport reuse),
- implement “sleep-based” rate limiting (jittery, wasteful).

This kata forces **explicit control** over *rate*, *in-flight concurrency*, and *cancellation*.

## 🎯 The Scenario
You’re building an internal service that needs to fetch user widgets from a downstream API:
- API allows **10 requests/sec** with bursts up to **20**
- Your service must also cap concurrency at **max 8 in-flight** requests
- If any request fails, cancel everything immediately (fail-fast), and return the first error.

## 🛠 The Challenge
Implement `FanOutClient` with:
- `FetchAll(ctx context.Context, userIDs []int) (map[int][]byte, error)`

### 1. Functional Requirements
- [ ] Requests must respect a **QPS rate limit** + **burst**.
- [ ] Requests must run concurrently but never exceed **MaxInFlight**.
- [ ] Results returned as `map[userID]payload`.
- [ ] On first error, cancel remaining work and return immediately.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must** use `golang.org/x/time/rate.Limiter`.
- [ ] **Must** use `golang.org/x/sync/semaphore.Weighted` (or equivalent semaphore pattern) for MaxInFlight.
- [ ] **Must** use `http.NewRequestWithContext`.
- [ ] **Must NOT** use `time.Sleep` for rate limiting.
- [ ] **Must** reuse a single `http.Client` (with a configured `Transport` + `Timeout`).
- [ ] Logging via `log/slog` (structured fields: userID, attempt, latency).

## 🧪 Self-Correction (Test Yourself)
- **If you spawn `len(userIDs)` goroutines:** you failed backpressure.
- **If cancellation doesn’t stop waiting callers:** you failed context propagation.
- **If QPS is enforced using `Sleep`:** you failed rate limiting.
- **If you use `http.DefaultClient`:** you failed HTTP hygiene.

## 📚 Resources
- https://pkg.go.dev/golang.org/x/time/rate
- https://pkg.go.dev/golang.org/x/sync/semaphore
- https://go.dev/src/net/http/client.go
- https://go.dev/src/net/http/transport.go
