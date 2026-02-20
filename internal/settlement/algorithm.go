package settlement

import (
	"container/heap"
	"github.com/user/fintech-expense-tracker/internal/model"
)

// UserBalance represents the net balance of a user in the group.
type UserBalance struct {
	UserID uint
	Name   string
	Amount int64
}

// BalanceHeap implements a Max-Heap for UserBalance, ensuring deterministic output.
type BalanceHeap []UserBalance

func (h BalanceHeap) Len() int           { return len(h) }
func (h BalanceHeap) Less(i, j int) bool {
	// Primary: Highest amount first (Max-Heap behavior)
	if h[i].Amount != h[j].Amount {
		return h[i].Amount > h[j].Amount
	}
	// Secondary: Tie-break with UserID ASC for bitwise-identical output across nodes
	return h[i].UserID < h[j].UserID
}
func (h BalanceHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *BalanceHeap) Push(x interface{}) { *h = append(*h, x.(UserBalance)) }
func (h *BalanceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// ComputeSettlements minimizes the number of transactions to resolve all debts.
// It uses a greedy approach by matching the largest creditor with the largest debtor.
// Time Complexity: O(N log N) where N is the number of participants.
func ComputeSettlements(netBalances []model.NetBalance) []model.Settlement {
	creditors := &BalanceHeap{}
	debtors := &BalanceHeap{}
	heap.Init(creditors)
	heap.Init(debtors)

	// Phase 1: Separate users into creditors and debtors
	for _, b := range netBalances {
		if b.Balance > 0 {
			heap.Push(creditors, UserBalance{UserID: b.UserID, Name: b.Name, Amount: b.Balance})
		} else if b.Balance < 0 {
			heap.Push(debtors, UserBalance{UserID: b.UserID, Name: b.Name, Amount: -b.Balance})
		}
	}

	var settlements []model.Settlement
	// Phase 2: Greedily match the largest creditor and debtor
	for creditors.Len() > 0 && debtors.Len() > 0 {
		creditor := heap.Pop(creditors).(UserBalance)
		debtor := heap.Pop(debtors).(UserBalance)

		// Settle the minimum of the two balances
		settleAmount := creditor.Amount
		if debtor.Amount < settleAmount {
			settleAmount = debtor.Amount
		}

		settlements = append(settlements, model.Settlement{
			FromUserID: debtor.UserID,
			FromUser:   debtor.Name,
			ToUserID:   creditor.UserID,
			ToUser:     creditor.Name,
			Amount:     settleAmount,
		})

		creditor.Amount -= settleAmount
		debtor.Amount -= settleAmount

		// Re-insert if there's remaining balance
		if creditor.Amount > 0 { heap.Push(creditors, creditor) }
		if debtor.Amount > 0 { heap.Push(debtors, debtor) }
	}
	return settlements
}
