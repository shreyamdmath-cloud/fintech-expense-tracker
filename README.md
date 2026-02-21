# Fintech-Grade Expense Settlement Engine

A high-performance, bitwise-deterministic expense settlement API built with Go. This engine is designed for financial precision, auditability, and massive scale.

## Key Features

### Deterministic Settlement Engine
- **O(N log N) Optimization**: Reduces complex debt chains into minimal transactions.
- **Secondary Sorting Guarantee**: Uses `(Amount DESC, UserID ASC)` for heap ordering, ensuring identical settlement results across distributed nodes.
- **Zero-Float Precision**: Uses `int64` (paise) for all calculations to prevent floating-point inaccuracies.

### Production Readiness
- **Panic Recovery**: Middleware to prevent server crashes on runtime exceptions.
- **Structured Logging**: Full request/response logging for audit trails.
- **Health Diagnostics**: Standardized `/health` endpoint for monitoring.
- **Centralized Errors**: Consistent API error responses for seamless integration.

### Intelligent Persistence
- **PostgreSQL Support**: Primary production database.
- **Zero-Config SQLite Fallback**: Automatically switches to local SQLite (`fintech_tracker.db`) if PostgreSQL is unavailable, perfect for rapid demos. Schema remains fully compatible with PostgreSQL.

## Financial Integrity & Precision

In financial systems, precision is non-negotiable. This engine employs several strategies to ensure absolute correctness:

- **The Float Trap**: Standard floating-point types (`float32`, `float64`) use binary fractions which cannot exactly represent base-10 decimals like 0.1 or 0.7. This leads to rounding errors that accumulate over thousands of transactions (rounding drift).
- **Smallest Currency Unit Modeling**: All monetary values are stored as `int64` representing the smallest unit (e.g., Paise for INR). Integer arithmetic is inherently exact and binary-stable.
- **Strict Sum Validation**: Custom validators ensure that the sum of discrete splits exactly equals the total transaction amount before persistence, preventing "missing cents" bugs.
- **Database Alignment**: We use `BIGINT` in PostgreSQL to align perfectly with Go's `int64`, ensuring no precision loss during I/O operations.

## Settlement Algorithm & Correctness

The settlement engine uses a greedy matching algorithm optimized for minimal transaction volume.

### Algorithm Properties:
- **Cycle Removal**: By converting pairwise debts into net balances, the system naturally eliminates debt cycles (A -> B -> C -> A), simplifying the graph into a bipartite sets of creditors and debtors.
- **Greedy matching**: Matching the largest creditor with the largest debtor at each step reduces the total number of transactions.
- **Efficiency**: The algorithm runs in $O(N \log N)$ time due to the use of Max-Heaps for tracking balances.
- **Termination**: The algorithm is guaranteed to terminate as each step resolves the full balance of at least one participant.
- **Transaction Bound**: Total transactions are guaranteed to be $\le (K - 1)$, where $K$ is the number of participants with non-zero balances.

### Deterministic Output
In distributed financial systems, reconciliation requires consistency. Our engine enforces **Bitwise Determinism**:
1. Net balances are computed for all users.
2. Users are sorted into Max-Heaps using `(Amount DESC, UserID ASC)`. 
3. Ties in debt amounts are broken by `UserID`, ensuring that the settlement sequence is identical regardless of the execution environment or node.

## Getting Started

### 1. Requirements
- **Go 1.23+**
- **Shell**: PowerShell (recommended) or Bash.

### 2. Run the System
```powershell
# Tidy dependencies
go mod tidy

# Run core logic tests
go test ./internal/settlement/... -v

# Start the API server
go run cmd/api/main.go
```

## Demonstration Audit
The following scenario was verified during the project audit:

1. **Group Created**: "Trip" (Alice & Bob).
2. **Expense**: Alice paid **1000** ($10.00), shared **500/500**.
3. **Settlement**: Engine correctly identified **Bob owes Alice 500**.

### Quick Verification Command:
```powershell
Invoke-RestMethod -Uri http://localhost:8080/health
```

## Architecture
- **Clean Architecture**: Decoupled Handler -> Service -> Repository layers.
- **Idempotency**: Support for `Idempotency-Key` to safely retry transactions.
- **Concurrent-Safe**: Designed for stateless horizontal scaling.

## Documentation
For more in-depth information about our engineering decisions and precision standards, please refer to:
- [Money Handling Approach](file:///C:/Users/Asus/.gemini/antigravity/scratch/fintech-expense-tracker/money_handling_approach.md): Detailed explanation of our "Zero-Float" policy and integer arithmetic.
- [AI Interaction Log](file:///C:/Users/Asus/.gemini/antigravity/scratch/fintech-expense-tracker/prompt.md): Audit trail of development prompts and architectural oversight.
