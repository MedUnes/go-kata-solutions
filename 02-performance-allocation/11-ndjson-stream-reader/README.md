# Kata 11: The NDJSON Reader That Survives Long Lines
**Target Idioms:** Streaming I/O (`io.Reader`), `bufio.Reader` vs `Scanner`, Handling `ErrBufferFull`, Low Allocation  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
Seasoned devs reach for `bufio.Scanner` and it “works”… until production sends a line > 64K and you get:
`bufio.Scanner: token too long`.

This kata forces you to implement a streaming reader that can handle **arbitrarily large lines** without falling over.

## 🎯 The Scenario
You ingest NDJSON logs from stdin or a file. Lines can be huge (hundreds of KB). You must process line-by-line.

## 🛠 The Challenge
Implement:
- `func ReadNDJSON(ctx context.Context, r io.Reader, handle func([]byte) error) error`

### 1. Functional Requirements
- [ ] Call `handle(line)` for each line (without the trailing newline).
- [ ] Stop immediately on `handle` error.
- [ ] Stop immediately on `ctx.Done()`.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must NOT** rely on default `bufio.Scanner` behavior.
- [ ] **Must** use `bufio.Reader` and correctly handle `ReadSlice('\n')` returning `ErrBufferFull`.
- [ ] **Must** avoid per-line allocations where possible (reuse buffers).
- [ ] Wrap errors with line number context using `%w`.

## 🧪 Self-Correction (Test Yourself)
- **If a 200KB line crashes with “token too long”:** you failed.
- **If cancellation doesn’t stop promptly:** you failed.
- **If you allocate a new buffer each line:** you failed the low-allocation goal.

## 📚 Resources
- https://pkg.go.dev/bufio
- https://pkg.go.dev/io
