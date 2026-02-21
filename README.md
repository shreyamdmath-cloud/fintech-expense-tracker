# Fintech-Grade Expense Settlement Engine

[![GitHub](https://img.shields.io/badge/GitHub-View%20Source-blue?logo=github)](https://github.com/shreyamdmath-cloud/fintech-expense-tracker)
🔗 **View Source on GitHub**: [shreyamdmath-cloud/fintech-expense-tracker](https://github.com/shreyamdmath-cloud/fintech-expense-tracker)

A high-performance, deterministic engine designed for financial precision, auditability, and minimal transaction settlement. This system reduces complex debt graphs into a minimal set of transactions while maintaining absolute correctness.

## 🧭 Documentation Links
[📘 Overview](index.html) | [⚙ Architecture](#system-architecture) | [💰 Money Handling](money_handling_approach.html) | [🧠 Settlement Algorithm](#5-settlement-algorithm-deep-dive) | [🔗 View on GitHub](https://github.com/shreyamdmath-cloud/fintech-expense-tracker)

## 1. High-Level Problem Context

Managing shared expenses (roommates, group trips, shared bills) often results in fragmented, inefficient debt chains. Naive tracking leads to redundant transactions (e.g., Alice pays Bob who pays Charlie). In large-scale financial systems, these "debt cycles" increase transaction costs and reconciliation overhead.

This engine solves the problem by transforming pairwise debts into a clean net-balance state and applying an optimized matching algorithm. Beyond simple arithmetic, it enforces **Bitwise Determinism**—ensuring that for any given set of inputs, every node in a distributed system arrives at the exact same settlement sequence.

## 2. System Architecture

The project follows the principles of **Clean Architecture**, enforcing a strict separation of concerns to ensure testability and maintainability.

### Workflow & Layering
```text
Client (HTTP Request)
       ↓
[Gin Web Router] (Middleware: Logger, Recovery)
       ↓
[Handlers] (Request Binding, Validation, API Response)
       ↓
[Services] (Domain Orchestration, Settlement Logic)
       ↓
[Repository Layer] (GORM Data Abstraction)
       ↓
[Database] (PostgreSQL Primary / SQLite Fallback)

Settlement Engine (internal package)
↳ Invoked by Service Layer
```

- **Gin handles routing**: Efficiently dispatches incoming HTTP requests to appropriate handlers.
- **Handlers manage request/response**: Handles JSON binding, validation, and lifecycle.
- **Services contain business logic**: Orchestrates domain-level operations and settlement computation.
- **Repository handles persistence**: Decouples domain logic from database-specific GORM implementations.
- **Settlement Engine is isolated**: Pure logic package for algorithm clarity and independent scalability.
- **Architecture**: Ensures strict separation of concerns and simplifies unit testing.

## 3. Data Model

The system uses a normalized relational model optimized for query performance and financial integrity.

- **Users**: Unique entities with name and email.
- **Groups**: Collaborative spaces for managing sets of participants.
- **Expenses**: Atomic financial events containing the total amount (as `BIGINT`) and an idempotency key.
- **ExpenseSplits**: A one-to-many relationship defining the exact share of an expense for each member of the group.
- **Precision Choice**: We use `BIGINT` in the database to align with Go's `int64`. This bypasses the inaccuracies of floating-point arithmetic at the infrastructure level.

## 4. API Documentation

### Endpoint Overview
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| POST | `/api/v1/users` | Create a new participant |
| POST | `/api/v1/groups` | Initialize a new group |
| POST | `/api/v1/groups/:id/members` | Add a user to a group |
| POST | `/api/v1/groups/:id/expenses` | Log an expense with splits |
| GET | `/api/v1/groups/:id/balances` | Query net balances for all members |
| GET | `/api/v1/groups/:id/settlements`| Calculate optimized transactions |

### Sample Integration (Add Expense)
**Request:** `POST /api/v1/groups/1/expenses`
```json
{
  "paid_by_id": 1,
  "description": "Team Dinner",
  "amount": 1000,
  "idempotency_key": "dinner-2024-01-01",
  "splits": {
    "1": 500,
    "2": 500
  }
}
```
**Response:** `201 Created`

## 5. Settlement Algorithm Deep Dive

The engine utilizes a greedy matching algorithm designed to resolve all debts in the minimum number of transactions.

1.  **Transformation**: All pairwise debts are flattened into a single "Net Balance" per user (e.g., if A owes B 10 and B owes A 5, the net is A: -5).
2.  **Cycle Elimination**: By moving to net balances, circular debts (A → B → C → A) are automatically eliminated.
3.  **Heuristic Matching**: The system populates two Max-Heaps: **Creditors** (positive balance) and **Debtors** (negative balance).
4.  **Greedy Resolve**: At each step, the largest creditor is matched with the largest debtor. This greedily reduces the remaining population of people with non-zero balances.
5.  **Deterministic Tie-Breaking**: If two users have identical balances, the system sorts by `UserID ASC`. This ensures the matching order is identical across execution nodes.
6.  **Termination**: The process repeats until all balances are zero. Total transactions are guaranteed to be $\le (K - 1)$, where $K$ is the number of participants.
7.  **Performance**: The usage of Max-Heaps ensures $O(N \log N)$ time complexity.

## 6. Deterministic Output

In distributed fintech systems, reconciliation depends on reproducibility. Naive greedy algorithms often fail here if they don't handle "ties" (e.g., two people owe the same amount).

Our engine enforces **Secondary Sorting**: (Amount DESC, UserID ASC). This guarantee ensures that if Alice and Bob both owe $50, the algorithm will always settle Alice first if her UserID is lower. Without this, reconciliation across distributed shards would be impossible.

## 7. Financial Integrity

### The Zero-Float Policy
Floating-point numbers (`float32`, `64`) use binary fractions which cannot exactly represent base-10 decimals. This leads to cumulative rounding drift.

- **Smallest Unit Modeling**: All values are handled as `int64` representing the smallest currency unit (e.g., Paise).
- **Conservation of Value**: Before persisting an expense, a custom validator ensures the sum of all splits exactly matches the total amount.
- **Integrity Boundary**: The API rejects any request where `sum(splits) != total_amount`.

## 8. Scalability & Production Considerations

- **Stateless Design**: The API is fully stateless, allowing for horizontal scaling behind a load balancer.
- **Indexing**: Database indices are placed on `group_id` across `expenses` and `splits` for sub-millisecond query performance.
- **Memory Complexity**: Settlement is performed in-memory per group, with $O(N)$ space complexity, making it highly efficient for thousands of concurrent group settlements.

## 9. Optimization Scenarios

### Example 1: Simple Settlement
**Before**: Alice paid 1000. Bob owes 500. Charlie owes 500.
**Result**:
- Bob → Alice: 500
- Charlie → Alice: 500

### Example 2: Circular Debt Elimination
**Debt Graph**: Alice owes Bob 100. Bob owes Charlie 100. Charlie owes Alice 100.
**Result**:
- **0 Transactions**. The engine identifies all net balances as 0.

## 10. Documentation
For further technical details:
- [Money Handling Approach](money_handling_approach.html)
- [AI Interaction Log](prompt.html)

---
*Fintech-Grade Settlement Engine • High Precision • High Performance*
