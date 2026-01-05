# Kata 04: The Interface-Based Middleware Chain
**Target Idioms:** Interface Design, Middleware Pattern, Composition over Inheritance
**Difficulty:** 🔴 Advanced

## 🧠 The "Why"
Object-oriented developers often reach for class hierarchies and inheritance when building pipelines. In Go, **interfaces enable composition over inheritance**. The unidiomatic approach is to create a `BaseHandler` class with virtual methods. The idiomatic Go way uses small interfaces composed together. This pattern powers `http.Handler`, `io.Reader`, and many standard library patterns - but developers from other ecosystems struggle to see when to split interfaces.

## 🎯 The Scenario
You're building a **real-time analytics pipeline** for user events. Each event must pass through multiple processing stages: validation, enrichment, filtering, and finally storage. New stages will be added frequently. The pipeline must be:
- Modular (add/remove stages without rewriting core logic)
- Observable (track metrics at each stage)
- Recoverable (continue processing after non-critical errors)

## 🛠 The Challenge
Create a middleware chain for processing user events.

### 1. Functional Requirements
* [ ] Process events through a configurable chain of middleware
* [ ] Each middleware can modify, filter, or reject events
* [ ] Provide metrics (counters, latencies) for each stage
* [ ] Support graceful shutdown with context cancellation

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
To pass this kata, you **must** strictly adhere to these rules:
* [ ] **Small Interfaces:** Define a `Processor` interface with a single method: `Process(context.Context, Event) ([]Event, error)`
* [ ] **Middleware Composition:** Each middleware must implement the `Processor` interface and wrap another `Processor`
* [ ] **Functional Options:** Configure middleware using functional options (e.g., `WithMetricsCollector()`)
* [ ] **Context Propagation:** All middleware must respect context cancellation
* [ ] **Zero Global State:** No package-level variables for configuration or state
* [ ] **Testable by Design:** Each middleware must be unit-testable in isolation

## 🧪 Self-Correction (Test Yourself)
Test your implementation against these scenarios:
1.  **The "Infinite Loop":**
    * Create a middleware that generates 2 events from 1 input
    * Chain it with a filtering middleware
    * **Fail Condition:** If events multiply uncontrollably or memory usage grows exponentially
2.  **The "Context Leak":**
    * Add a middleware with a 10s timeout
    * Cancel the context after 1s
    * **Fail Condition:** If any middleware continues processing after context cancellation
3.  **The "Interface Pollution":**
    * Try to add a new middleware that needs access to database connections
    * **Fail Condition:** If you had to modify the core `Processor` interface to add database methods

## 📚 Resources
* [Go Proverbs by Rob Pike](https://go-proverbs.github.io/)
* [The Go Blog: Lexical Scanning in Go](https://blog.golang.org/lexical-scanning)
* [Standard Library Inspiration: net/http.Handler](https://pkg.go.dev/net/http#Handler)
* [Small Interfaces in the Standard Library](https://medium.com/@cep21/small-interfaces-in-go-1e912a7a7883)