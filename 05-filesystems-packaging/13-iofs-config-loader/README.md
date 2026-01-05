# Kata 13: The Filesystem-Agnostic Config Loader
**Target Idioms:** `io/fs` abstraction, `fs.WalkDir`, Testability via `fstest.MapFS`, `embed` readiness  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
In Go, passing `"/etc/app/config"` all over the place hard-couples your logic to the OS.
Idiomatic Go uses `fs.FS` so you can:
- load from disk,
- load from embedded files,
- load from a ZIP filesystem,
- unit test without touching the real filesystem.

## 🎯 The Scenario
Your CLI loads configuration fragments from a directory tree, merges them, and prints a final config report.

## 🛠 The Challenge
Implement:
- `func LoadConfigs(fsys fs.FS, root string) (map[string][]byte, error)`

### 1. Functional Requirements
- [ ] Walk `root` recursively and read all `*.conf` files.
- [ ] Return a map of `path -> content`.
- [ ] Reject invalid paths cleanly.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Must** accept `fs.FS` (not `os` paths) in the core API.
- [ ] **Must** use `fs.WalkDir` and `fs.ReadFile`.
- [ ] **Must NOT** use `os.Open` / `filepath.Walk` inside the core loader.
- [ ] Unit tests must use `testing/fstest.MapFS`.

## 🧪 Self-Correction (Test Yourself)
- **If you can’t test without real files:** you failed.
- **If your loader only works on disk:** you failed the abstraction goal.

## 📚 Resources
- https://pkg.go.dev/io/fs
- https://go.dev/src/embed/embed.go
