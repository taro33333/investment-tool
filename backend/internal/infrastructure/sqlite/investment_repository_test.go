package sqlite

import (
	"context"
	"database/sql"
	"moneyget/internal/domain"
	"testing"
)

func TestInvestmentRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewInvestmentRepository(db)
	ctx := context.Background()

	// テストデータの作成
	money, err := domain.NewMoney(1000, "JPY")
	if err != nil {
		t.Fatalf("Failed to create money: %v", err)
	}

	investment, err := domain.NewInvestment(
		domain.NewInvestmentID("test-investment"),
		money,
		domain.Stock,
		domain.Conservative,
	)
	if err != nil {
		t.Fatalf("Failed to create investment: %v", err)
	}

	// Create のテスト
	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, investment)
		if err != nil {
			t.Errorf("Failed to create investment: %v", err)
		}
	})

	// FindByID のテスト
	t.Run("FindByID", func(t *testing.T) {
		found, err := repo.FindByID(ctx, investment.ID())
		if err != nil {
			t.Errorf("Failed to find investment: %v", err)
		}
		if found.ID().Value != investment.ID().Value {
			t.Errorf("Expected ID %s, got %s", investment.ID().Value, found.ID().Value)
		}
		if found.Amount().Amount != investment.Amount().Amount {
			t.Errorf("Expected amount %f, got %f", investment.Amount().Amount, found.Amount().Amount)
		}
		if found.Type() != investment.Type() {
			t.Errorf("Expected type %s, got %s", investment.Type(), found.Type())
		}
	})

	// Save のテスト
	t.Run("Save", func(t *testing.T) {
		newMoney, _ := domain.NewMoney(2000, "JPY")
		investment.UpdateAmount(newMoney)
		err := repo.Save(ctx, investment)
		if err != nil {
			t.Errorf("Failed to save investment: %v", err)
		}

		// 更新された値を確認
		found, err := repo.FindByID(ctx, investment.ID())
		if err != nil {
			t.Errorf("Failed to find investment after update: %v", err)
		}
		if found.Amount().Amount != 2000 {
			t.Errorf("Expected updated amount 2000, got %f", found.Amount().Amount)
		}
	})

	// FindAll のテスト
	t.Run("FindAll", func(t *testing.T) {
		investments, err := repo.FindAll(ctx)
		if err != nil {
			t.Errorf("Failed to find all investments: %v", err)
		}
		if len(investments) != 1 {
			t.Errorf("Expected 1 investment, got %d", len(investments))
		}
	})

	// Delete のテスト
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, investment.ID())
		if err != nil {
			t.Errorf("Failed to delete investment: %v", err)
		}

		// 削除されたことを確認
		_, err = repo.FindByID(ctx, investment.ID())
		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, got %v", err)
		}
	})
}
