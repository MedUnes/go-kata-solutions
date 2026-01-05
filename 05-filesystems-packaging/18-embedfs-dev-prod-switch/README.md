# Kata 18: embed.FS Dev/Prod Switch Without Handler Forks
**Target Idioms:** `embed`, `io/fs`, Build Tags, `fs.Sub`, Same Handler Code Path  
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
Embedding assets is great for production (single binary), but terrible for frontend iteration if every CSS tweak needs a rebuild.
Idiomatic Go solves this with:
- compile-time selection via build tags
- a shared `fs.FS` abstraction so handler code doesn’t branch on “dev/prod”.

## 🎯 The Scenario
You run a small internal dashboard:
- Prod: ship a single binary (assets embedded).
- Dev: designers update `static/` and `templates/` live without recompiling.

## 🛠 The Challenge
Create a server that serves:
- templates from `templates/`
- static assets from `static/`

### 1. Functional Requirements
- [ ] `GET /` renders an HTML template.
- [ ] `GET /static/...` serves static files.
- [ ] Dev mode serves from disk; prod mode serves embedded.
- [ ] Handler code is identical in both modes.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
- [ ] **Build tags:** two files:
    - `assets_dev.go` with `//go:build dev`
    - `assets_prod.go` with `//go:build !dev`
- [ ] **Return `fs.FS`:** `func Assets() (templates fs.FS, static fs.FS, err error)`
- [ ] **Use `fs.Sub`:** exported FS must have *clean roots* (no `static/static/...` path bugs).
- [ ] **No runtime env checks in handlers:** mode selection must be compile-time.
- [ ] **Single `http.FileServer` setup:** no duplicated handler logic for dev vs prod.

## 🧪 Self-Correction (Test Yourself)
1. **Live Reload**
    - Build with `-tags dev`.
    - Modify a CSS file and refresh.
    - **Pass:** change shows without rebuild.

2. **Binary Portability**
    - Build without tags.
    - Delete `static/` and `templates/` from disk.
    - **Pass:** server still serves assets/templates.

3. **Prefix Correctness**
    - Request `/static/app.css`.
    - **Pass:** works in both modes (no 404 due to prefix mismatch).

## 📚 Resources
- https://pkg.go.dev/embed
- https://pkg.go.dev/io/fs
- https://pkg.go.dev/io/fs#Sub
