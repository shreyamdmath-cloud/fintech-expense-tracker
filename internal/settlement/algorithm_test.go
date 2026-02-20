package settlement

import (
	"reflect"
	"testing"
	"github.com/user/fintech-expense-tracker/internal/model"
)

func TestComputeSettlements(t *testing.T) {
	tests := []struct {
		name string
		nb []model.NetBalance
		expectedLen int
	}{
		{"Simple", []model.NetBalance{{UserID:1, Name:"A", Balance:500}, {UserID:2, Name:"B", Balance:-500}}, 1},
		{"Chain", []model.NetBalance{{UserID:1, Name:"A", Balance:1000}, {UserID:2, Name:"B", Balance:-400}, {UserID:3, Name:"C", Balance:-600}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ComputeSettlements(tt.nb)
			if len(s) != tt.expectedLen { t.Errorf("expected %d, got %d", tt.expectedLen, len(s)) }
		})
	}
}

func TestDeterminism(t *testing.T) {
	// Case with equal balances to trigger secondary sorting rule
	nb := []model.NetBalance{
		{UserID: 2, Name: "B", Balance: 1000},
		{UserID: 1, Name: "A", Balance: 1000},
		{UserID: 3, Name: "C", Balance: -2000},
	}
	s1 := ComputeSettlements(nb)
	
	// Expected: C pays A first because UserID 1 < UserID 2
	if s1[0].ToUserID != 1 {
		t.Errorf("determinism failed: expected UserID 1 to be settled first due to tie-breaker, got %d", s1[0].ToUserID)
	}

	s2 := ComputeSettlements(nb)
	if !reflect.DeepEqual(s1, s2) { 
		t.Errorf("not deterministic") 
	}
}
