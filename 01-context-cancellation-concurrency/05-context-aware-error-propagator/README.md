# Kata 05: The Context-Aware Error Propagator
**Target Idioms:** Error Wrapping, Context-Aware Errors, Custom Error Types
**Difficulty:** 🟡 Intermediate

## 🧠 The "Why"
Developers from dynamic languages often treat errors as simple strings. Java developers wrap exceptions in layers of inheritance. **Go's error philosophy is different:** errors are values that should carry context and be inspectable without string parsing. The unidiomatic pattern is to `log.Printf("error: %v", err)` and return nil - this destroys debugging context. Idiomatic Go preserves the original error while adding layers of context.

## 🎯 The Scenario
You're building a **cloud storage gateway** that interacts with multiple services: authentication, metadata database, and blob storage. When a file upload fails, operators need to know exactly which layer failed and why - was it auth timeout? database deadlock? storage quota exceeded? Your error handling must preserve this information while being safe for logging.

## 🛠 The Challenge
Create a service that uploads files to cloud storage with proper error handling.

### 1. Functional Requirements
* [ ] Implement three layers: `AuthService`, `MetadataService`, `StorageService`
* [ ] Each layer can fail with specific error types
* [ ] Return errors that expose the failure point and original cause

### 2. The "Idiomatic" Constraints (Pass/Fail Criteria)
To pass this kata, you **must** strictly adhere to these rules:
* [ ] **NO string-based error inspection:** You must use `%w` with `fmt.Errorf` for wrapping
* [ ] **Custom Error Types:** Create specific error types for each service layer (e.g., `AuthError`, `StorageQuotaError`)
* [ ] **Context-Aware Errors:** Errors must implement `Timeout()` and `Temporary()` methods where appropriate
* [ ] **Safe Logging:** Errors must redact sensitive information (API keys, credentials) when logged
* [ ] **Error Unwrapping:** Your errors must support `errors.Is()` and `errors.As()` for programmatic inspection

## 🧪 Self-Correction (Test Yourself)
Test your error handling with these scenarios:
1.  **The "Sensitive Data Leak":**
    * Force an auth error with a mock API key
    * **Fail Condition:** If `fmt.Sprint(err)` contains the API key string
2.  **The "Lost Context":**
    * Wrap an `AuthError` three times through different layers
    * **Fail Condition:** If `errors.As(err, &AuthError{})` returns false
3.  **The "Timeout Confusion":**
    * Create a timeout error in the storage layer
    * **Fail Condition:** If `errors.Is(err, context.DeadlineExceeded)` returns false

## 📚 Resources
* [Go 1.13 Error Wrapping](https://go.dev/blog/go1.13-errors)
* [Error Handling in Upspin](https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html)
* [Don't just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
