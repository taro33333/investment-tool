package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"moneyget/internal/domain"
	"os"
)

type transactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) domain.TransactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// パニックの場合は元のパニックを優先し、ロールバックエラーはログに記録するなどの処理を検討
				fmt.Printf("Rollback error after panic: %v\n", rbErr)
			}
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %v)", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}

func (tm *transactionManager) WithTransaction(ctx context.Context) (*sql.Tx, error) {
	return tm.db.BeginTx(ctx, nil)
}

func RunMigrations(db *sql.DB) error {
	// スキーマファイルを読み込み
	schema, err := os.ReadFile("internal/infrastructure/sqlite/schema.sql")
	if err != nil {
		return err
	}

	// トランザクション内でマイグレーションを実行
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(schema))
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v (original error: %v)", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}
