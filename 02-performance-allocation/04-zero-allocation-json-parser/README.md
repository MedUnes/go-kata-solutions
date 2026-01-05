# Kata 04: The Zero-Allocation JSON Parser

**Target Idioms:** Performance Optimization, `json.RawMessage`, Streaming Parsers, Buffer Reuse
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why**
Developers from dynamic languages often parse JSON by unmarshaling entire documents into `map[string]interface{}` or generic structs. In high-throughput Go services, this creates:
1. Massive memory churn (GC pressure)
2. Unnecessary allocations for unused fields
3. Lost type safety

The Go way: **Parse only what you need, reuse everything**. This kata teaches you to treat JSON as a stream, not a document.

## 🎯 The Scenario
You're processing **10MB/s of IoT sensor data** with JSON like:
```json
{"sensor_id": "temp-1", "timestamp": 1234567890, "readings": [22.1, 22.3, 22.0], "metadata": {...}}
```
You only need `sensor_id` and the first reading value. Traditional unmarshal would allocate for all fields and the entire readings array.

## 🛠 The Challenge
Implement `SensorParser` that extracts specific fields without full unmarshaling.

### 1. Functional Requirements
* [ ] Parse `sensor_id` (string) and first `readings` value (float64) from JSON stream
* [ ] Process `io.Reader` input (could be HTTP body, file, or network stream)
* [ ] Handle malformed JSON gracefully (skip bad records, continue parsing)
* [ ] Benchmark under 100ns per object and 0 allocations per parse

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
* [ ] **NO `encoding/json.Unmarshal`**: Use `json.Decoder` with `Token()` streaming
* [ ] **Reuse Buffers**: Use `sync.Pool` for `bytes.Buffer` or `json.Decoder`
* [ ] **Early Exit**: Stop parsing once required fields are found
* [ ] **Type Safety**: Return concrete struct `SensorData{sensorID string, value float64}`, not `interface{}`
* [ ] **Memory Limit**: Process arbitrarily large streams in constant memory (<1MB heap)

## 🧪 Self-Correction (Test Yourself)
1. **The Allocation Test**:
   ```go
   go test -bench=. -benchmem -count=5
   ```
   **Pass**: `allocs/op` = 0 for parsing loop
   **Fail**: Any allocations in hot path

2. **The Stream Test**:
    - Pipe 1GB of JSON through your parser (mock with repeating data)
    - **Pass**: Memory usage flatlines after warm-up
    - **Fail**: Memory grows linearly with input size

3. **The Corruption Test**:
    - Input: `{"sensor_id": "a"} {"bad json here` (malformed second object)
    - **Pass**: Returns first object, logs/skips second, doesn't panic
    - **Fail**: Parser crashes or stops processing entirely

## 📚 Resources
* [Go JSON Stream Parsing](https://ahmet.im/blog/golang-json-stream-parse/)
* [json.RawMessage Tutorial](https://www.sohamkamani.com/golang/json/#raw-messages)
* [Advanced JSON Techniques](https://eli.thegreenplace.net/2019/go-json-cookbook/)