# Money Handling Approach: Financial Precision in Go

## 1. The Core Problem: Floating Point Inaccuracy
In financial systems, using floating-point types (like `float32` or `float64`) is a fundamental error. Due to how IEEE 754 represents decimals in binary, operations like `0.1 + 0.2` often result in `0.30000000000000004`. This "drift" is unacceptable in bill-splitting applications where every paisa/cent must be accounted for.

## 2. Our Solution: The "Zero-Float" Policy
Our system employs a **Zero-Float Policy** throughout the entire stack (API, Domain Logic, and Database).

### A. Smallest Unit Representation
We store all monetary values as **Integers** (`int64`). 
- **Unit**: The smallest currency unit (e.g., Paise for INR, Cents for USD).
- **Example**: An expense of ₹10.50 is stored and processed as `1050`.

### B. Scalability & Limits
Using `int64` allows us to handle amounts up to $9 \times 10^{18}$ units. For INR, this is roughly ₹90 quadrillion, far exceeding the needs of any group expense application.

## 3. Database Integrity
In the PostgreSQL/SQLite schema, we use the `BIGINT` type. 
- **Decimal Types**: While SQL `DECIMAL` or `NUMERIC` types are "proper," they often require specialized library handling in application code (like `shopspring/decimal`). 
- **Decision**: By using `BIGINT` at the DB level and `int64` in Go, we achieve the highest performance and simplest code without sacrificing even a single unit of precision.

## 4. Split Validation (The Conservation Law)
To maintain financial integrity, the system enforces a strict "Conservation of Value" check:
```go
totalSplits := 0
for _, share := range splits {
    totalSplits += share
}
if totalSplits != expenseAmount {
    return Error("Sum of splits must exactly equal the total amount")
}
```
This ensures that no "hanging paise" are created during the splitting process.

## 5. Settlement Precision
The $O(N \log N)$ settlement algorithm operates exclusively on these integer balances. Because the sum of all net balances in any group is guaranteed to be exactly `0`, the algorithm will always settle perfectly without any fractional remainders.

---
[← Back to Overview](file:///C:/Users/Asus/.gemini/antigravity/scratch/fintech-expense-tracker/README.md) | [AI Interaction Log →](file:///C:/Users/Asus/.gemini/antigravity/scratch/fintech-expense-tracker/prompt.md)
