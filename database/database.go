package database

import (
	"context"
	"fmt"
	"time"

	"wb/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type DB struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*DB, error) {
	log.Info().Msg("connect to DB")

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	db := &DB{pool: pool}
	
	if err = db.simpleMigrate(ctx); err != nil {
		return nil, fmt.Errorf("error on migrations: %w", err)
	}

	return db, nil
}

type TransactionOperator struct {
	txPool pgx.Tx
}

func (db *DB) BeginTransaction(ctx context.Context) (*TransactionOperator, error) {
	txPool, err := db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, err
	}

	return &TransactionOperator{txPool: txPool}, nil
}

// todo
func (db *DB) simpleMigrate(ctx context.Context) error {
	log.Info().Msg("run migration - init single table")

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	sql := `CREATE TABLE IF NOT EXISTS employee (
				id UUID NOT NULL PRIMARY KEY,
				fullname VARCHAR(255) NOT NULL DEFAULT '',
				gender INTEGER NOT NULL DEFAULT 0,
				age INTEGER NOT NULL DEFAULT 0,
				email VARCHAR(50) NOT NULL DEFAULT '',
				address TEXT NOT NULL DEFAULT '',
				vacation_days INTEGER NOT NULL DEFAULT 0,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
				deleted_at TIMESTAMP
			)`

	if _, err = tx.txPool.Exec(ctx, sql); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (tr *TransactionOperator) Commit(ctx context.Context) error {
	return tr.txPool.Commit(ctx)
}

func (tr *TransactionOperator) Rollback(ctx context.Context) error {
	return tr.txPool.Rollback(ctx)
}

func (tr *TransactionOperator) InsertEmployee(ctx context.Context, employee domain.Employee) error {
	sql := `INSERT INTO employee (id, fullname, gender, age, email, address) 
				VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := tr.txPool.Exec(
		ctx, sql,
		employee.ID, employee.Fullname, employee.Gender,
		employee.Age, employee.Email, employee.Address,
	); err != nil {
		return err
	}

	return nil
}

func (tr *TransactionOperator) DeleteEmployee(ctx context.Context, employeeID uuid.UUID) error {
	sql := `UPDATE employee SET deleted_at=$2 WHERE id=$1 AND deleted_at IS NULL`

	if _, err := tr.txPool.Exec(ctx, sql, employeeID, time.Now()); err != nil {
		return err
	}

	return nil
}

func (tr *TransactionOperator) GetVacationDays(ctx context.Context, employeeID uuid.UUID) (uint, error) {
	sql := `SELECT vacation_days FROM employee WHERE id=$1 AND deleted_at IS NULL`

	var result uint

	if err := tr.txPool.QueryRow(ctx, sql, employeeID).Scan(&result); err != nil {
		return 0, err
	}

	return result, nil
}

func (tr *TransactionOperator) SearchEmployee(ctx context.Context, query string) ([]*domain.Employee, error) {
	sql := `SELECT id, fullname, gender, age, email, address, created_at 
			FROM employee
			WHERE LOWER(fullname) LIKE $1 || '%'`

	rows, err := tr.txPool.Query(ctx, sql, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*domain.Employee, 0)

	for rows.Next() {
		employee := new(domain.Employee)

		if err = rows.Scan(
			&employee.ID,
			&employee.Fullname,
			&employee.Gender,
			&employee.Age,
			&employee.Email,
			&employee.Address,
			&employee.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, employee)
	}

	return result, rows.Err()
}
