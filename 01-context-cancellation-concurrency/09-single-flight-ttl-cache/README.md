# Kata 09: The Cache Stampede Shield
**Target Idioms:** `singleflight`, TTL Cache, DoChan + Context Select, Lock Avoidance  
**Difficulty:** 🔴 Advanced

## 🧠 The "Why"
In many stacks, caching is “just Redis”. In Go, an in-process cache is common, but people:
- hold locks while calling the loader (deadly),
- refresh the same key N times concurrently (stampede),
- can’t cancel waiters cleanly.

This kata is about **deduplicating in-flight loads** and making waiters **context-cancellable**.

## 🎯 The Scenario
You have expensive per-key loads (e.g., DB or remote API). If 200 goroutines ask for the same key at once:
- loader must run **once**
- others must wait (or return on ctx cancel)
- TTL must be enforced

## 🛠 The Challenge
Implement:
- `type Cache[K comparable, V any] struct { ... }`
- `Get(ctx context.Context, key K, loader func(context.Context) (V, error)) (V, error)`

### 1. Functional Requirements
- [ ] Return cached value if not expired.
- [ ] If expired/missing: load once, share result to all callers.
- [ ] Callers must be able to stop waiting via `ctx.Done()`.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must** use `golang.org/x/sync/singleflight.Group`.
- [ ] **Must** use `DoChan` + `select` on `ctx.Done()` to cancel waiters.
- [ ] **Must NOT** hold a mutex while calling `loader`.
- [ ] Errors must be wrapped with key context using `%w`.

## 🧪 Self-Correction (Test Yourself)
- **If 200 goroutines trigger 200 loads:** you failed (no stampede protection).
- **If a canceled context still blocks waiting:** you failed.
- **If you lock around loader execution:** you failed (contention / deadlocks).

## 📚 Resources
- https://pkg.go.dev/golang.org/x/sync/singleflight
- https://go.dev/blog/go1.13-errors
