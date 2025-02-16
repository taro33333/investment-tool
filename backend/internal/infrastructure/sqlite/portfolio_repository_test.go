package sqlite

import (
	"context"
	"database/sql"
	"moneyget/internal/domain"
	"testing"
)

func TestPortfolioRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewPortfolioRepository(db)
	investmentRepo := NewInvestmentRepository(db)
	ctx := context.Background()

	// テストデータの準備
	portfolio := domain.NewPortfolio(
		domain.NewPortfolioID("test-portfolio"),
		"test-user",
	)

	// 投資データの準備
	money, _ := domain.NewMoney(1000, "JPY")
	investment, _ := domain.NewInvestment(
		domain.NewInvestmentID("test-investment"),
		money,
		domain.InvestmentTypeStock,
		domain.InvestmentStrategyValue,
	)
	
	// 投資を作成
	err := investmentRepo.Create(ctx, investment)
	if err != nil {
		t.Fatalf("Failed to create investment: %v", err)
	}

	// ポートフォリオに投資を追加
	err = portfolio.AddInvestment(investment)
	if err != nil {
		t.Fatalf("Failed to add investment to portfolio: %v", err)
	}

	// Create のテスト
	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, portfolio)
		if err != nil {
			t.Errorf("Failed to create portfolio: %v", err)
		}
	})

	// FindByID のテスト
	t.Run("FindByID", func(t *testing.T) {
		found, err := repo.FindByID(ctx, portfolio.ID())
		if err != nil {
			t.Errorf("Failed to find portfolio: %v", err)
		}
		if found.ID().Value != portfolio.ID().Value {
			t.Errorf("Expected ID %s, got %s", portfolio.ID().Value, found.ID().Value)
		}
		if found.UserID != portfolio.UserID {
			t.Errorf("Expected UserID %s, got %s", portfolio.UserID, found.UserID)
		}
		if len(found.GetInvestments()) != 1 {
			t.Errorf("Expected 1 investment, got %d", len(found.GetInvestments()))
		}
	})

	// FindByUserID のテスト
	t.Run("FindByUserID", func(t *testing.T) {
		found, err := repo.FindByUserID(ctx, portfolio.UserID)
		if err != nil {
			t.Errorf("Failed to find portfolio by user ID: %v", err)
		}
		if found.ID().Value != portfolio.ID().Value {
			t.Errorf("Expected ID %s, got %s", portfolio.ID().Value, found.ID().Value)
		}
	})

	// FindByInvestmentID のテスト
	t.Run("FindByInvestmentID", func(t *testing.T) {
		found, err := repo.FindByInvestmentID(ctx, investment.ID())
		if err != nil {
			t.Errorf("Failed to find portfolio by investment ID: %v", err)
		}
		if found.ID().Value != portfolio.ID().Value {
			t.Errorf("Expected ID %s, got %s", portfolio.ID().Value, found.ID().Value)
		}
	})

	// Update のテスト
	t.Run("Update", func(t *testing.T) {
		portfolio.UserID = "updated-user"
		err := repo.Update(ctx, portfolio)
		if err != nil {
			t.Errorf("Failed to update portfolio: %v", err)
		}

		found, err := repo.FindByID(ctx, portfolio.ID())
		if err != nil {
			t.Errorf("Failed to find updated portfolio: %v", err)
		}
		if found.UserID != "updated-user" {
			t.Errorf("Expected updated UserID 'updated-user', got %s", found.UserID)
		}
	})

	// Delete のテスト
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, portfolio.ID())
		if err != nil {
			t.Errorf("Failed to delete portfolio: %v", err)
		}

		// 削除されたことを確認
		_, err = repo.FindByID(ctx, portfolio.ID())
		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, got %v", err)
		}
	})
}