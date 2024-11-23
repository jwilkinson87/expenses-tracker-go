package repository

import (
	"context"

	"example.com/expenses-tracker/internal/models"
)

type ExpenseRepository interface {
	CreateExpense(context.Context, *models.Expense) error
	GetExpense(context.Context, string) (*models.Expense, error)
	UpdateExpense(context.Context, *models.Expense) error
	DeleteExpense(context.Context, *models.Expense) error
	GetAllForUser(context.Context, *models.User) (models.Expenses, error)
}

type UserRepository interface {
	CreateUser(context.Context, *models.User) error
	UpdateUser(context.Context, *models.User) error
	GetUserByEmailAddress(context.Context, string) (*models.User, error)
	DeleteUser(context.Context, *models.User) error
}

type UserAuthRepository interface {
	CreateAuthToken(context.Context, *models.UserToken) error
	DeleteAuthToken(context.Context, *models.UserToken) error
	GetByAuthToken(context.Context, string) (*models.UserToken, error)
}
