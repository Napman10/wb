package service

import (
	"context"

	"github.com/google/uuid"

	"wb/database"
	"wb/domain"
)

type Service struct {
	db *database.DB
}

func New(db *database.DB) IService {
	return &Service{db: db}
}

func (s *Service) HireEmployee(ctx context.Context, employee domain.Employee) (uuid.UUID, error) {
	tx, err := s.db.BeginTransaction(ctx)
	if err != nil {
		return uuid.Nil, nil
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	employee.ID = uuid.New()

	if err = tx.InsertEmployee(ctx, employee); err != nil {
		return uuid.Nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return employee.ID, nil
}

func (s *Service) FireEmployee(ctx context.Context, employeeID uuid.UUID) error {
	tx, err := s.db.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	if err = tx.DeleteEmployee(ctx, employeeID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) GetVacationDays(ctx context.Context, employeeID uuid.UUID) (uint, error) {
	tx, err := s.db.BeginTransaction(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	result, err := tx.GetVacationDays(ctx, employeeID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return result, nil
}

func (s *Service) SearchEmployee(ctx context.Context, query string) ([]*domain.Employee, error) {
	tx, err := s.db.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	result, err := tx.SearchEmployee(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return result, nil
}
