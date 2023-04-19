package database

import (
	"context"

	"github.com/google/uuid"

	"wb/domain"
)

type TransactionManager interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	InsertEmployee(ctx context.Context, employee domain.Employee) error
	DeleteEmployee(ctx context.Context, employeeID uuid.UUID) error
	GetVacationDays(ctx context.Context, employeeID uuid.UUID) (uint, error)
	SearchEmployee(ctx context.Context, query string) (*domain.Employee, error)
}
