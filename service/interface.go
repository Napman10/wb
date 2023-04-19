package service

import (
	"context"

	"github.com/google/uuid"

	"wb/domain"
)

type IService interface {
	HireEmployee(ctx context.Context, employee domain.Employee) (uuid.UUID, error)
	FireEmployee(ctx context.Context, employeeID uuid.UUID) error
	GetVacationDays(ctx context.Context, employeeID uuid.UUID) (uint, error)
	SearchEmployee(ctx context.Context, query string) ([]*domain.Employee, error)
}
