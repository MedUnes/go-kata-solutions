# Kata 08: The Retry Policy That Respects Context
**Target Idioms:** Retry Classification, Error Wrapping (`%w`), Timer Reuse, Context Deadlines  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
In other languages, retries are often hidden in SDKs. In Go, it’s easy to write:
- infinite retry loops,
- retry-on-any-error (bad),
- retry that ignores context cancellation (worse),
- retry implemented with repeated `time.Sleep` (hard to test, wasteful).

This kata makes you implement a **testable**, **context-aware** retry loop.

## 🎯 The Scenario
You call a flaky downstream service. You should retry only on **transient** failures:
- `net.Error` with `Timeout() == true`
- HTTP 429 / 503 (if you model HTTP)
- sentinel `ErrTransient`

Everything else must fail immediately.

## 🛠 The Challenge
Implement:
- `type Retryer struct { ... }`
- `func (r *Retryer) Do(ctx context.Context, fn func(context.Context) error) error`

### 1. Functional Requirements
- [ ] Retries up to `MaxAttempts`.
- [ ] Uses exponential backoff: `base * 2^attempt` with a max cap.
- [ ] Optional jitter (deterministic in tests).
- [ ] Stops immediately on `ctx.Done()`.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must NOT** call `time.Sleep` inside the retry loop.
- [ ] **Must** use a `time.Timer` and `Reset` it (timer reuse).
- [ ] **Must** wrap the final error with context (attempt count) using `%w`.
- [ ] **Must** classify errors using `errors.Is` / `errors.As`.

## 🧪 Self-Correction (Test Yourself)
- **If context cancellation only stops after the sleep:** you failed.
- **If you retry non-transient errors:** you failed classification.
- **If you can’t test it without real time:** inject time/jitter sources.

## 📚 Resources
- https://go.dev/blog/go1.13-errors
- https://pkg.go.dev/errors
- https://pkg.go.dev/time
