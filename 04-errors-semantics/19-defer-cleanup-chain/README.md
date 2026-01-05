# Kata 19: The Cleanup Chain (defer + LIFO + Error Preservation)
**Target Idioms:** `defer` Discipline, Named Returns, Error Composition (`errors.Join`), Close/Rollback Ordering  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
`defer` is easy to misuse:
- deferring in loops (resource spikes),
- ignoring `Close()` / `Rollback()` errors,
- losing the original failure when cleanup also fails,
- wrong cleanup ordering (commit then rollback nonsense).

Idiomatic Go keeps cleanup local, ordered, and preserves important errors.

## 🎯 The Scenario
You implement `BackupDatabase`:
- open output file
- connect DB
- begin transaction
- stream rows to file
- commit
  If anything fails, you must close/rollback what was already acquired.

## 🛠 The Challenge
Implement:
- `func BackupDatabase(ctx context.Context, dbURL, filename string) (err error)`

Use mock interfaces for DB + Tx + Rows if you want (recommended).

### 1. Functional Requirements
- [ ] Open file for writing.
- [ ] Connect to DB.
- [ ] Begin Tx.
- [ ] Write data (simulate streaming).
- [ ] Commit on success.
- [ ] On failure: rollback + close resources in correct order.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Defer cleanup immediately after acquisition.**
- [ ] **No manual cleanup paths** except by controlling flags (e.g., `committed bool`) used by deferred funcs.
- [ ] **Preserve both errors:** if main operation fails and cleanup fails too, return a combined error (`errors.Join`).
- [ ] **Named return `err`** so defers can amend it safely.
- [ ] **No defer-in-loop for per-row resources:** if your mock has per-row closers, show the correct pattern.

## 🧪 Self-Correction (Test Yourself)
1. **Tx Begin Fails**
    - Make `Begin()` error.
    - **Pass:** file + db connection still close.

2. **Commit Fails + Close Fails**
    - Make `Commit()` return error and also make `file.Close()` return error.
    - **Pass:** returned error clearly contains both (use `errors.Join`).

3. **No FD Leak**
    - Run 1000 times.
    - **Pass:** file descriptors don’t grow.

## 📚 Resources
- https://go.dev/blog/defer-panic-and-recover
- https://go.dev/doc/go1.20 (errors.Join)
