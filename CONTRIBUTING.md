# Contributing to Go Katas

First off, thank you for your interest in **Go Katas**.

This repository is a curated collection of **constraint-driven, production-grade Go challenges**.

---

## ⚠️ Scope & Intent (Please Read)

**This project is curated by a developer **transitioning to Go**, for developers transitioning to Go.**

**I value seniority, experience, and collective knowledge, and believe it is more than ever crucial to us, in the era of "AI" and information overflow**

The value of this repo comes from:

- encoding real **production failure modes**
- making **idiomatic Go constraints explicit**
- encouraging deliberate, corrective practice

👉 **If you are a seasoned Go engineer and spot an incorrect or dangerous pattern, please open an Issue or PR.**
Expert review and correction is not just welcome, it is essential.

This repository aims to converge toward correctness through community scrutiny.

---

## 🛠 How to Use This Repo (For Learners)

1. **Fork** this repository.
2. **Solve the katas in your own fork or locally.**
3. **Do NOT submit Pull Requests with personal solution implementations.**
    - This repo is for *practice prompts*, not solution sharing.
4. Optionally compare your work against:
    - curated **reference implementations** (where provided)
    - your own past solutions

---

## 🚀 How to Contribute (For Experienced Engineers)

We accept Pull Requests in the following categories:

### 1. New Katas (Highly Encouraged)

If you have encountered a real production issue, mismatch, or failure mode:

- Place yourself under the corresponding category folder (numbered 01-XXX to 06-XXX under the root project folder). If you think it is worth a new category, please suggest it
- Create a new folder following the `XX-topic-name` convention
- Copy the [README_TEMPLATE](README_TEMPLATE.md) file into the new folder's root and rename it to `README.md`
- Edit the `README.md`, for instance:
    - **The Why** (what breaks in production)
    - **The Scenario**
    - **Idiomatic Constraints** (pass/fail criteria)
- Focus on **software engineering concerns**, such as:
    - concurrency & cancellation
    - memory and allocation behavior
    - HTTP and I/O hygiene
    - observability and lifecycle management

❌ Avoid:

- toy algorithm puzzles
- framework-specific abstractions
- "clever" tricks without production relevance

---

### 2. Reference Implementations (Curated)

Some katas include **reference implementations**.

Reference implementations:

- exist to **demonstrate constraints**, not personal style
- must be minimal, readable, and defensible
- should include tests or benchmarks where appropriate
- may be challenged or replaced if better patterns exist

If you believe a reference solution is:

- inefficient
- unsafe in production
- or not idiomatic

👉 **Please submit a PR with rationale and references.**

---

### 3. Fixing Explanations & Constraints

We welcome PRs that:

- clarify the "Why" section of a kata
- correct incorrect assumptions
- improve wording, accuracy, or structure
- add missing edge cases or constraints

Well-argued Issues are just as valuable as code.

---

## 🧭 Governance & Decision-Making

- Decisions prioritize:
    - production safety
    - clarity over cleverness
    - testability over opinion
    - alignment with Go’s standard library and documented idioms

When tradeoffs exist, they should be **documented explicitly**, not hidden.

---

## 🧠 Philosophy

These katas are designed to help experienced developers:

- unlearn habits from other ecosystems
- develop Go-specific instincts
- reason about failure, not just success

If you disagree with a constraint, **that is a feature, not a bug**.
Please challenge it constructively.

---

Thank you for contributing to serious, deliberate Go practice.
