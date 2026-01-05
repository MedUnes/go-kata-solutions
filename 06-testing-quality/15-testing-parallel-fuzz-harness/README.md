# Kata 15: The Go Test Harness (Subtests, Parallel, Fuzz)
**Target Idioms:** Table-Driven Tests, `t.Run`, `t.Parallel`, Fuzzing (`go test -fuzz`)  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
Developers often write:
- one-off tests with repetition,
- unsafe parallel subtests (loop variable capture),
- no fuzz testing for parsers/sanitizers.

Idiomatic Go testing is:
- table-driven,
- readable failures,
- parallel where safe,
- fuzzed for edge cases.

## 🎯 The Scenario
You’re implementing a sanitizer:
- `func NormalizeHeaderKey(s string) (string, error)`
  Rules:
- only ASCII letters/digits/hyphen allowed
- normalize to canonical header form (e.g., `content-type` -> `Content-Type`)
- reject invalid input

## 🛠 The Challenge
Write:
1) The implementation, and
2) A test suite that proves it’s solid.

### 1. Functional Requirements
- [ ] Canonicalize valid inputs.
- [ ] Reject invalid characters.
- [ ] Stable behavior (same input => same output).

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] Tests must be **table-driven** with `t.Run`.
- [ ] Use **parallel subtests** correctly (no loop var capture bugs).
- [ ] Include a **fuzz test** that:
    - never panics,
    - never returns a string containing invalid characters,
    - roundtrips canonical form (calling Normalize twice is idempotent).

## 🧪 Self-Correction (Test Yourself)
- **If parallel subtests flake:** you likely captured the loop variable.
- **If fuzzing finds panics:** you missed an edge case.

## 📚 Resources
- https://go.dev/blog/subtests
- https://go.dev/wiki/TableDrivenTests
- https://go.dev/doc/security/fuzz/
- https://go.dev/doc/tutorial/fuzz
