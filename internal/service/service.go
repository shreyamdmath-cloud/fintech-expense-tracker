package service

import (
	"errors"
	"github.com/user/fintech-expense-tracker/internal/model"
	"github.com/user/fintech-expense-tracker/internal/repository"
	"github.com/user/fintech-expense-tracker/internal/settlement"
)

type Service interface {
	CreateUser(name, email string) (*model.User, error)
	CreateGroup(name string) (*model.Group, error)
	AddUserToGroup(groupID, userID uint) error
	AddExpense(groupID, paidByID uint, description string, amount int64, splits map[uint]int64, idempotencyKey string) (*model.Expense, error)
	GetBalances(groupID uint) ([]model.NetBalance, error)
	GetSettlements(groupID uint) ([]model.Settlement, error)
}

type expenseService struct {
	repo repository.Repository
}

func NewExpenseService(repo repository.Repository) Service {
	return &expenseService{repo: repo}
}

func (s *expenseService) CreateUser(name, email string) (*model.User, error) {
	user := &model.User{Name: name, Email: email}
	err := s.repo.CreateUser(user)
	return user, err
}

func (s *expenseService) CreateGroup(name string) (*model.Group, error) {
	group := &model.Group{Name: name}
	err := s.repo.CreateGroup(group)
	return group, err
}

func (s *expenseService) AddUserToGroup(groupID, userID uint) error {
	return s.repo.AddUserToGroup(groupID, userID)
}

func (s *expenseService) AddExpense(groupID, paidByID uint, description string, amount int64, splits map[uint]int64, idempotencyKey string) (*model.Expense, error) {
	if amount <= 0 { return nil, errors.New("INVALID_AMOUNT: amount must be > 0") }
	if len(splits) == 0 { return nil, errors.New("EMPTY_SPLITS: expense must have at least one split") }
	var totalShare int64
	var expenseSplits []model.ExpenseSplit
	for userID, share := range splits {
		if share < 0 { return nil, errors.New("NEGATIVE_SHARE: individual shares cannot be negative") }
		totalShare += share
		expenseSplits = append(expenseSplits, model.ExpenseSplit{UserID: userID, Share: share})
	}
	if totalShare != amount { return nil, errors.New("SPLIT_MISMATCH: sum of shares does not equal total amount") }
	expense := &model.Expense{
		GroupID: groupID, PaidByID: paidByID, Description: description,
		Amount: amount, IdempotencyKey: idempotencyKey, Splits: expenseSplits,
	}
	err := s.repo.CreateExpense(expense)
	return expense, err
}

func (s *expenseService) GetBalances(groupID uint) ([]model.NetBalance, error) {
	group, err := s.repo.GetGroupByID(groupID)
	if err != nil { return nil, err }
	if len(group.Members) == 0 { return nil, errors.New("EMPTY_GROUP: group has no members") }
	expenses, err := s.repo.GetExpensesByGroupID(groupID)
	if err != nil { return nil, err }

	balances := make(map[uint]int64)
	userNames := make(map[uint]string)
	for _, m := range group.Members {
		balances[m.ID] = 0
		userNames[m.ID] = m.Name
	}

	for _, e := range expenses {
		balances[e.PaidByID] += e.Amount
		for _, sp := range e.Splits {
			balances[sp.UserID] -= sp.Share
		}
	}

	var netBalances []model.NetBalance
	for id, b := range balances {
		if b != 0 {
			netBalances = append(netBalances, model.NetBalance{UserID: id, Name: userNames[id], Balance: b})
		}
	}
	return netBalances, nil
}

func (s *expenseService) GetSettlements(groupID uint) ([]model.Settlement, error) {
	nb, err := s.GetBalances(groupID)
	if err != nil { return nil, err }
	return settlement.ComputeSettlements(nb), nil
}
