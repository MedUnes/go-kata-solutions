# Kata 10: The Worker Pool With Backpressure and Joined Errors
**Target Idioms:** Worker Pools, Channel Ownership, `errors.Join`, Context Cancellation  
**Difficulty:** 🔴 Advanced

## 🧠 The "Why"
Many devs bring “thread pool” instincts and end up with:
- goroutines that never exit,
- unbounded queues,
- “first error wins” even when you want a summary,
- ad-hoc error channels without cleanup.

This kata forces correctness: **bounded work**, **clean shutdown**, and **error aggregation**.

## 🎯 The Scenario
You process a stream of jobs (e.g., image resizing). You want:
- fixed number of workers
- bounded queue (backpressure)
- either fail-fast OR collect all errors (configurable)

## 🛠 The Challenge
Implement:
- `type Pool struct { ... }`
- `Run(ctx context.Context, jobs <-chan Job) error`

Where `Job` is `func(context.Context) error`.

### 1. Functional Requirements
- [ ] `N` workers process from `jobs`.
- [ ] Optional `StopOnFirstError`.
- [ ] If not fail-fast: return `errors.Join(errs...)` after draining.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must** use `errors.Join` for aggregation.
- [ ] **Must** respect `ctx.Done()` (workers exit).
- [ ] **Must** close internal channels from the sender side only.
- [ ] **Must** guarantee no goroutine leak when `jobs` closes early or ctx cancels.

## 🧪 Self-Correction (Test Yourself)
- **If workers keep running after ctx cancel:** failed.
- **If you can deadlock by closing channels from the wrong side:** failed.
- **If you return before draining in non-fail-fast mode:** failed.

## 📚 Resources
- https://go.dev/doc/go1.20
- https://go.dev/src/errors/join.go
