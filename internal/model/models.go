package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Groups    []*Group  `gorm:"many2many:group_members;" json:"groups,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Group struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Members   []*User   `gorm:"many2many:group_members;" json:"members"`
	Expenses  []Expense `json:"expenses,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Expense struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	GroupID        uint           `gorm:"index:idx_group_id;not null" json:"group_id"`
	PaidByID       uint           `gorm:"index:idx_paid_by;not null" json:"paid_by_id"`
	PaidBy         User           `gorm:"foreignKey:PaidByID" json:"paid_by"`
	Description    string         `gorm:"size:500;not null" json:"description"`
	Amount         int64          `gorm:"type:bigint;not null" json:"amount"`
	IdempotencyKey string         `gorm:"size:100;uniqueIndex;not null" json:"idempotency_key"`
	Splits         []ExpenseSplit `json:"splits"`
	CreatedAt      time.Time      `json:"created_at"`
}

type ExpenseSplit struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	ExpenseID uint  `gorm:"index:idx_expense_id;not null" json:"expense_id"`
	UserID    uint  `gorm:"index:idx_user_id;not null" json:"user_id"`
	User      User  `gorm:"foreignKey:UserID" json:"user"`
	Share     int64 `gorm:"type:bigint;not null" json:"share"`
}

type NetBalance struct {
	UserID  uint   `json:"user_id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

type Settlement struct {
	FromUserID uint   `json:"from_user_id"`
	FromUser   string `json:"from_user"`
	ToUserID   uint   `json:"to_user_id"`
	ToUser     string `json:"to_user"`
	Amount     int64  `json:"amount"`
}
