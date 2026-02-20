package repository

import (
	"github.com/user/fintech-expense-tracker/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	CreateGroup(group *model.Group) error
	GetGroupByID(id uint) (*model.Group, error)
	AddUserToGroup(groupID, userID uint) error
	CreateExpense(expense *model.Expense) error
	GetExpensesByGroupID(groupID uint) ([]model.Expense, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateUser(user *model.User) error { return r.db.Create(user).Error }
func (r *gormRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}
func (r *gormRepository) CreateGroup(group *model.Group) error { return r.db.Create(group).Error }
func (r *gormRepository) GetGroupByID(id uint) (*model.Group, error) {
	var group model.Group
	err := r.db.Preload("Members").First(&group, id).Error
	return &group, err
}
func (r *gormRepository) AddUserToGroup(groupID, userID uint) error {
	var group model.Group
	if err := r.db.First(&group, groupID).Error; err != nil { return err }
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil { return err }
	return r.db.Model(&group).Association("Members").Append(&user)
}
func (r *gormRepository) CreateExpense(expense *model.Expense) error {
	return r.db.Create(expense).Error
}
func (r *gormRepository) GetExpensesByGroupID(groupID uint) ([]model.Expense, error) {
	var expenses []model.Expense
	err := r.db.Preload("Splits").Where("group_id = ?", groupID).Find(&expenses).Error
	return expenses, err
}
