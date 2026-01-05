# Kata 17: The Context-Aware Channel Sender (No Leaked Producers)
**Target Idioms:** Pipeline Cancellation, Select-on-Send, Channel Ownership, Goroutine Leak Prevention  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
A goroutine sending on a channel blocks until a receiver is ready (unless buffered space is available). If the receiver exits early (timeout, HTTP cancel, upstream error), producers can block forever and leak.

Idiomatic Go fixes this by:
- threading `context.Context` through the pipeline
- **selecting on every send** (`case out <- v` vs `case <-ctx.Done()`), as recommended in Go’s pipeline cancellation patterns and real leak writeups.

## 🎯 The Scenario
You’re building a data pipeline step that fetches N URLs concurrently and streams results downstream. If the request is canceled (client disconnect, global timeout), **all fetchers must stop immediately** and no goroutine may remain blocked on `out <- result`.

## 🛠 The Challenge
Implement:
- `type DataFetcher struct { ... }`
- `func (f *DataFetcher) Fetch(ctx context.Context, urls []string) <-chan Result`

Where:
- `Result` contains `URL`, `Body []byte`, `Err error` (or similar).

### 1. Functional Requirements
- [ ] Start concurrent fetchers for all URLs (or bounded concurrency if you choose).
- [ ] Send results as they complete (order doesn’t matter).
- [ ] Stop promptly on `ctx.Done()`.
- [ ] Close the output channel exactly once after all producers exit.
- [ ] Return partial results that already completed before cancellation.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Every send uses select:** no bare `out <- x`.
- [ ] **Channel ownership:** only the producer side closes `out`.
- [ ] **No goroutine leaks:** all goroutines exit when ctx is canceled.
- [ ] **No double close:** prove it structurally (single closer goroutine). Avoid `sync.Once` unless you can justify it.
- [ ] **Buffer choice is intentional:** if you buffer, document why and how you chose the size.

### 3. Hints (Allowed Tools)
- You may use `errgroup` or a simple worker pattern, but the key is: **send must be cancel-aware**.
- If you do bounded concurrency, prefer `x/sync/semaphore` or a worker pool (but don’t turn this kata into a rate-limiter kata).

## 🧪 Self-Correction (Test Yourself)
1. **Forgotten Sender**
    - Start 50 fetchers, consume only 1 result, then cancel.
    - **Pass:** goroutine count returns near baseline quickly (use `runtime.NumGoroutine()` as a sanity check).

2. **Cancellation Before First Receive**
    - Cancel ctx immediately after calling `Fetch`.
    - **Pass:** no goroutine blocks trying to send.

3. **Close Discipline**
    - Cancel ctx from multiple places.
    - **Pass:** no `panic: close of closed channel`.

## 📚 Resources
- https://go.dev/blog/pipelines
- https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html
