# Kata 14: The Leak-Free Scheduler
**Target Idioms:** `time.Timer`/`time.Ticker`, Stop/Reset patterns, Jitter, Context Cancellation  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
Scheduling in Go is deceptively easy until you ship:
- goroutines that never stop,
- overlapping executions,
- ticker drift and backlog,
- resource retention from careless timer usage.

This kata makes you build a scheduler that is **predictable** and **stoppable**.

## 🎯 The Scenario
You need to periodically refresh a local cache:
- every 5s, with ±10% jitter
- do not overlap refreshes
- stop immediately on shutdown

## 🛠 The Challenge
Implement:
- `type Scheduler struct { ... }`
- `func (s *Scheduler) Run(ctx context.Context, job func(context.Context) error) error`

### 1. Functional Requirements
- [ ] Run `job` periodically (interval + jitter).
- [ ] Never run `job` concurrently with itself.
- [ ] Exit on `ctx.Done()`.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must NOT** use `time.Tick` (no stop control).
- [ ] **Must** use a `time.Timer` or `time.Ticker` with correct stop/reset.
- [ ] **Must** propagate context into `job`.
- [ ] Log job duration and errors via `slog`.

## 🧪 Self-Correction (Test Yourself)
- **If `job` overlap occurs:** you failed.
- **If cancel doesn’t stop quickly:** you failed.
- **If goroutines remain after exit:** you failed.

## 📚 Resources
- https://pkg.go.dev/time
- https://go.dev/wiki/Go123Timer
