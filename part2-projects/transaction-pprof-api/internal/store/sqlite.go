package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Account struct {
	ID      int `json:"id"`
	Balance int `json:"balance"`
}

type SQLiteStore struct {
	DB *sql.DB
}

func OpenSQLite(dsn string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	s := &SQLiteStore{DB: db}
	if err := s.initSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *SQLiteStore) initSchema(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS accounts (
  id INTEGER PRIMARY KEY,
  balance INTEGER NOT NULL
);
`)
	if err != nil {
		return err
	}
	// seed 两个账户（幂等）
	_, _ = s.DB.ExecContext(ctx, `INSERT OR IGNORE INTO accounts(id, balance) VALUES (1, 1000), (2, 1000);`)
	return nil
}

func (s *SQLiteStore) ListAccounts(ctx context.Context) ([]Account, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT id, balance FROM accounts ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Account
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Balance); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Transfer 事务示例：from 扣款、to 加款；任一步失败回滚
func (s *SQLiteStore) Transfer(ctx context.Context, fromID, toID, amount int) error {
	if amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if fromID == toID {
		return errors.New("from and to must differ")
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var fromBalance int
	if err := tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = ?`, fromID).Scan(&fromBalance); err != nil {
		return fmt.Errorf("load from account: %w", err)
	}
	if fromBalance < amount {
		return errors.New("insufficient balance")
	}

	if _, err := tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - ? WHERE id = ?`, amount, fromID); err != nil {
		return fmt.Errorf("debit: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + ? WHERE id = ?`, amount, toID); err != nil {
		return fmt.Errorf("credit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

