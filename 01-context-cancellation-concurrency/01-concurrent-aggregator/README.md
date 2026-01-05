# Kata 01: The Fail-Fast Data Aggregator

**Target Idioms:** Concurrency Control (`errgroup`), Context Propagation, Functional Options
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
In other languages, you might use `Promise.all` or strict thread pools to fetch data in parallel. In Go, seasoned developers often start with `sync.WaitGroup`, but quickly realize it lacks two critical features for production: **Error Propagation** and **Context Cancellation**.

If you spawn 10 goroutines and the first one fails, `WaitGroup` blindly waits for the other 9 to finish. **Idiomatic Go fails fast.**

## 🎯 The Scenario
You are building a **User Dashboard Backend**. To render the dashboard, you must fetch data from two independent, mock microservices:
1.  **Profile Service** (Returns "Name: Alice")
2.  **Order Service** (Returns "Orders: 5")

You need to fetch these in parallel to reduce latency. However, if *either* fails, or if the global timeout is reached, the entire operation must abort immediately to save resources.

## 🛠 The Challenge
Create a `UserAggregator` struct and a method `Aggregate(id int)` that orchestrates this fetching.

### 1. Functional Requirements
* [ ] The aggregator must be configurable (timeout, logger) without a massive constructor.
* [ ] Both services must be queried concurrently.
* [ ] The result should combine both outputs: `"User: Alice | Orders: 5"`.

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
To pass this kata, you **must** strictly adhere to these rules:

* [ ] **NO `sync.WaitGroup`:** You must use `golang.org/x/sync/errgroup`.
* [ ] **NO "Parameter Soup":** You must use the **Functional Options Pattern** for the constructor (e.g., `New(WithTimeout(2s))`).
* [ ] **Context is King:** You must pass `context.Context` as the first argument to your methods.
* [ ] **Cleanup:** If the Profile service fails, the Order service request must be cancelled (via Context) immediately.
* [ ] **Modern Logging:** Use `log/slog` for structured logging.

## 🧪 Self-Correction (Test Yourself)
Run your code against these edge cases:

1.  **The "Slow Poke":** * Set your aggregator timeout to `1s`.
    * Mock one service to take `2s`.
    * **Pass Condition:** Does your function return `context deadline exceeded` after exactly 1s?
2.  **The "Domino Effect":**
    * Mock the Profile Service to return an error immediately.
    * Mock the Order Service to take 10 seconds.
    * **Pass Condition:** Does your function return the error *immediately*? (If it waits 10s, you failed context cancellation).

## 📚 Resources
* [Go Concurrency: errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)
* [Functional Options for Friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)