# Kata 03: The Graceful Shutdown Server
**Target Idioms:** Context Propagation, Signal Handling, Channel Coordination, Resource Cleanup
**Difficulty:** 🔴 Advanced

## 🧠 The "Why"
In other ecosystems, graceful shutdown is often framework magic (Spring's `@PreDestroy`, Django's `close()`). Go forces explicit lifecycle management. The mismatch: developers used to automatic cleanup often leak goroutines, drop in-flight requests, or corrupt data during shutdown.

The Go way: **Own your lifecycle**. Every goroutine you spawn must have a controlled shutdown path.

## 🎯 The Scenario
Build an **HTTP Server with Background Worker** that must:
1. Accept HTTP requests (handled by a pool of worker goroutines)
2. Run a background cache warmer every 30 seconds
3. Maintain persistent database connections
4. Shutdown within 10 seconds when receiving SIGTERM, completing in-flight requests but rejecting new ones

## 🛠 The Challenge
Implement `Server` struct with `Start() error` and `Stop(ctx context.Context) error` methods.

### 1. Functional Requirements
* [ ] HTTP server on configurable port with request timeout
* [ ] Worker pool (configurable size) processes requests via channel
* [ ] Background cache warmer ticks every 30s (use `time.Ticker`)
* [ ] Database connection pool (mock with `net.Conn`)
* [ ] SIGTERM/SIGINT triggers graceful shutdown
* [ ] Shutdown completes within deadline or forces exit

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
* [ ] **Single Context Tree**: Root `context.Context` passed to `Start()`, canceled on shutdown
* [ ] **Channel Coordination**: Use `chan struct{}` for worker pool shutdown, not boolean flags
* [ ] **Proper Ticker Cleanup**: `defer ticker.Stop()` with select in goroutine
* [ ] **Dependency Order**: Shutdown in reverse order (stop accepting → drain workers → stop warmer → close DB)
* [ ] **No `os.Exit()` in business logic**: Shutdown should be testable without process termination

## 🧪 Self-Correction (Test Yourself)
1. **The Sudden Death Test**:
    - Send 100 requests, immediately send SIGTERM
    - **Pass**: Server completes in-flight requests (not all 100), logs "shutting down", closes cleanly
    - **Fail**: Server accepts new requests after signal, leaks goroutines, or crashes

2. **The Slow Leak Test**:
    - Run server for 5 minutes with 1 request/second
    - Send SIGTERM, wait 15 seconds
    - **Pass**: `go test` shows no goroutine leaks (use `runtime.NumGoroutine()`)
    - **Fail**: Any increase in goroutine count from start to finish

3. **The Timeout Test**:
    - Start long-running request (sleep 20s)
    - Send SIGTERM with 5s timeout context
    - **Pass**: Forces shutdown after 5s, logs "shutdown timeout"
    - **Fail**: Waits full 20s or deadlocks

## 📚 Resources
* [Go Blog: Context](https://go.dev/blog/context)
* [Graceful Shutdown in Go](https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a)
* [Signal Handling](https://medium.com/@marcus.olsson/writing-a-go-app-with-graceful-shutdown-5de1d2c6de96)
