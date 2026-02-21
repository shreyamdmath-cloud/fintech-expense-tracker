# AI Interaction Log: Fintech Expense Tracker

This document provides a transparent audit trail of the AI prompts used during the development of the high-precision settlement engine. It demonstrates responsible AI usage and architectural oversight.

## Responsible AI Usage
All AI-generated code was subject to human-in-the-loop review, rigorous unit testing, and manual audit for financial correctness. AI was utilized as a force multiplier for boilerplate generation and algorithm optimization, while the core business logic (deterministic tie-breaking and integer-based money handling) was strictly enforced by engineering requirements.

## Build Prompts
Used for scaffolding the core architecture and database layers.

### 1. Engine & Algorithm Specification
> "Write a Go package for a deterministic bill settlement algorithm. Use dual priority queues (Max-Heaps) to greedily match debtors and creditors in O(N log N) time. Crucially, implement a tie-breaking rule using UserID ASC to ensure bitwise-identical output across distributed nodes. Use `int64` representing Paise to prevent floating-point drift."

### 2. Database & GORM Layer
> "Design a normalized database schema for an expense tracker. Tables: `Users`, `Groups`, `Expenses` (with `idempotency_key`), and `ExpenseSplits`. In GORM, use `BIGINT` for all amount fields to match the integer-safe policy. Ensure GORM automatically creates associated splits when an expense is saved."

### 3. Clean Architecture API 
> "Implement a REST API using the Gin framework. Follow a strict layered architecture: Handlers (JSON binding/I/O), Service (Domain orchestration), and Repository (GORM data access). The DB initialization must gracefully fallback to a local SQLite `fintech_tracker.db` if the PostgreSQL DSN is unreachable."

## Audit & Refinement Prompts
Used for hardening, bug fixing, and ensuring production readiness.

### 4. Production Middleware & Error Consistency
> "Add production middleware to the Gin engine: Recovery from panics and Structured Logging. Implement a centralized `RespondWithError` helper that returns a standardized JSON structure: `{code, message, error}`. Create a `/health` diagnostic endpoint that pings the database."

### 5. Bug Fixing & Refinement
> "Fix a GORM constraint violation in `CreateExpense` where splits are being double-inserted. Refactor the handler `AddUserToGroup` to correctly use `c.ShouldBindJSON` to capture the `user_id` from the request body."

---
[← Back to Overview](README.md) | [Money Handling Approach →](money_handling_approach.md)
