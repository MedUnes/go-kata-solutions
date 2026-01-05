# Kata 02: The Concurrent Map with Sharded Locks

**Target Idioms:** Concurrency Safety, Map Sharding, `sync.RWMutex`, Avoiding `sync.Map` Pitfalls
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
Seasoned developers coming from Java might reach for `ConcurrentHashMap`-style solutions, while Pythonistas might think of GIL-protected dictionaries. In Go, you have three main options:
1. **Naive sync.Mutex around a map** (bottlenecks under high concurrency)
2. **sync.Map** (optimized for specific "append-only, read-heavy" cases, but opaque and often misused)
3. **Sharded maps** (manual control, maximized throughput)

The Go way is explicit control: if you know your access patterns, build a solution that fits. This kata forces you to understand *when* and *why* to choose sharding over sync.Map.

## 🎯 The Scenario
You're building a real-time **API Rate Limiter** that tracks request counts per user ID. The system handles 50k+ RPS with 95% reads (checking limits) and 5% writes (incrementing counters). A single mutex would serialize all operations-unacceptable. `sync.Map` might work but obscures memory usage and lacks type safety.

## 🛠 The Challenge
Implement `ShardedMap[K comparable, V any]` with configurable shard count that provides safe concurrent access.

### 1. Functional Requirements
* [ ] Type-safe generic implementation (Go 1.18+)
* [ ] `Get(key K) (V, bool)` - returns value and existence flag
* [ ] `Set(key K, value V)` - inserts or updates
* [ ] `Delete(key K)` - removes key
* [ ] `Keys() []K` - returns all keys (order doesn't matter)
* [ ] Configurable number of shards at construction

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
* [ ] **NO `sync.Map`**: Implement sharding manually with `[]map[K]V` and `[]sync.RWMutex`
* [ ] **Smart Sharding**: Use `fnv64` hashing for key distribution (don't rely on Go's random map iteration)
* [ ] **Read Optimization**: Use `RLock()` for `Get()` operations when safe
* [ ] **Zero Allocation Hot-Path**: `Get()` and `Set()` must not allocate memory in the critical section (no string conversion, no boxing)
* [ ] **Clean `Keys()`**: Implement without data races, even while concurrent writes occur

## 🧪 Self-Correction (Test Yourself)
1. **The Contention Test**:
    - Run 8 goroutines doing only `Set()` operations with sequential keys
    - With 1 shard: Should see heavy contention (use `go test -bench=. -cpuprofile` to verify)
    - With 64 shards: Should see near-linear scaling

2. **The Memory Test**:
    - Store 1 million `int` keys with `interface{}` values
    - **Fail Condition**: If your solution uses more than 50MB extra memory vs baseline map
    - **Hint**: Avoid `string(key)` conversions; use type-safe hashing

3. **The Race Test**:
    - Run `go test -race` with concurrent read/write/delete operations
    - Any race condition = automatic failure

## 📚 Resources
* [Go Maps Don't Appear to be O(1)](https://dave.cheney.net/2018/05/29/how-the-go-runtime-implements-maps-efficiently-without-generics)
* [When to use sync.Map](https://dave.cheney.net/2017/07/30/should-i-use-sync-map)
* [Practical Sharded Maps](https://github.com/orcaman/concurrent-map)
