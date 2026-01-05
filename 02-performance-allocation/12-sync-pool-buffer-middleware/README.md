# Kata 12: The sync.Pool Buffer Middleware
**Target Idioms:** `sync.Pool`, Avoiding GC Pressure, `bytes.Buffer` Reset, Benchmarks (`-benchmem`)  
**Difficulty:** 🔴 Advanced

## 🧠 The "Why"
In Go, performance regressions often come from allocation/GC churn, not “slow CPU”.
People use `sync.Pool` incorrectly:
- pooling long-lived objects (wrong),
- forgetting to reset buffers (data leak),
- storing huge buffers back into the pool (memory bloat).

This kata is about **safe pooling** for high-throughput handlers.

## 🎯 The Scenario
You’re writing an HTTP middleware that:
- reads up to 16KB of request body for audit logging
- must not allocate per-request in the hot path

## 🛠 The Challenge
Implement a middleware:
- `func AuditBody(max int, next http.Handler) http.Handler`

### 1. Functional Requirements
- [ ] Read up to `max` bytes of request body (do not consume beyond `max`).
- [ ] Log the captured bytes with `slog` fields.
- [ ] Pass the request downstream intact (body still readable).

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must** use `sync.Pool` to reuse buffers.
- [ ] **Must** `Reset()`/clear buffers before putting back.
- [ ] **Must** bound memory: never keep buffers larger than `max` in the pool.
- [ ] Provide a benchmark showing reduced allocations (`go test -bench . -benchmem`).

## 🧪 Self-Correction (Test Yourself)
- **If a request leaks previous request content:** you failed (no reset).
- **If allocations are ~O(requests):** you failed pooling.
- **If buffers grow unbounded and stay in pool:** you failed memory bounds.

## 📚 Resources
- https://pkg.go.dev/sync
- https://go.dev/doc/gc-guide
- https://go.dev/blog/pprof
