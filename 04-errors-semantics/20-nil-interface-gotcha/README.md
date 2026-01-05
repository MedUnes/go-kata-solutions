# Kata 20: The “nil != nil” Interface Trap (Typed nil Errors)
**Target Idioms:** Interface Semantics, Typed nil Pitfall, Safe Error Returns, `errors.As`  
**Difficulty:** 🔴 Advanced

## 🧠 The "Why"
In Go, an interface value is only nil when **both** its dynamic type and value are nil.
If you return a **typed nil pointer** (e.g., `(*MyError)(nil)`) as an `error`, the interface has a non-nil type, so `err != nil` becomes true even though the pointer inside is nil.

This bites real code in production (especially custom error types and factories).

## 🎯 The Scenario
A function returns `error`. Sometimes it returns a typed nil pointer.
Your caller checks `if err != nil` and takes an error path, logs misleading failures, or even panics when accessing fields/methods.

## 🛠 The Challenge
Write a minimal package that:
1) demonstrates the bug, and
2) fixes it with an idiomatic pattern.

### 1. Functional Requirements
- [ ] Implement `type MyError struct { Op string }` (or similar).
- [ ] Implement a function `DoThing(...) error` that **sometimes returns** `(*MyError)(nil)` as `error`.
- [ ] Demonstrate:
    - `err != nil` is true
    - `fmt.Printf("%T %#v\n", err, err)` shows the typed nil behavior
- [ ] Provide a corrected version that returns a true nil interface when there is no error.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must show the failing behavior** in a test (`go test`).
- [ ] **Must show the fix** in a test.
- [ ] **Must not “fix” by panicking or by sentinel errors.**
- [ ] Use one of these idiomatic fixes:
    - return `nil` explicitly when the pointer is nil
    - or return `error(nil)` in the relevant branch
- [ ] Demonstrate safe extraction using:
    - `var me *MyError; errors.As(err, &me)` and check `me != nil`

## 🧪 Self-Correction (Test Yourself)
1. **The Trap Repro**
    - Make `DoThing()` return `var e *MyError = nil; return e`
    - **Pass:** your test proves `err != nil` is true.

2. **The Fix**
    - If internal pointer is nil, return literal `nil`.
    - **Pass:** `err == nil` works, callers behave correctly.

3. **Extraction Safety**
    - Wrap the error and still extract with `errors.As`.
    - **Pass:** extraction works through wrapping layers.

## 📚 Resources
- https://go.dev/blog/laws-of-reflection (interface basics)
- https://go.dev/blog/go1.13-errors (errors.As)
- https://forum.golangbridge.org/t/logic-behind-failing-nil-check/16331
