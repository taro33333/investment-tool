package usecase

import (
	"context"
	"fmt"
	"moneyget/internal/domain"
	"moneyget/internal/domain/service"
	"testing"
)

func TestPortfolioUseCase_CreatePortfolio(t *testing.T) {
	ctx := context.Background()
	portfolioRepo := newMockPortfolioRepository()
	txManager := &mockTransactionManager{}
	eventPublisher := &mockEventPublisher{}
	strategyService := service.NewInvestmentStrategyService()

	useCase := NewPortfolioUseCase(
		portfolioRepo,
		txManager,
		eventPublisher,
		strategyService,
	)

	tests := []struct {
		name        string
		userID      string
		expectError bool
	}{
		{
			name:        "create portfolio successfully",
			userID:      "test-user",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portfolio, err := useCase.CreatePortfolio(ctx, tt.userID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if portfolio == nil {
					t.Error("Expected portfolio, got nil")
				}
				if portfolio.UserID != tt.userID {
					t.Errorf("Expected user ID %s, got %s", tt.userID, portfolio.UserID)
				}
				if len(portfolio.GetInvestments()) != 0 {
					t.Error("New portfolio should have no investments")
				}
			}
		})
	}
}

func TestPortfolioUseCase_GetPortfolioAnalysis(t *testing.T) {
	ctx := context.Background()
	portfolioRepo := newMockPortfolioRepository()
	txManager := &mockTransactionManager{}
	eventPublisher := &mockEventPublisher{}
	strategyService := service.NewInvestmentStrategyService()

	useCase := NewPortfolioUseCase(
		portfolioRepo,
		txManager,
		eventPublisher,
		strategyService,
	)

	// テスト用のポートフォリオを作成
	portfolio := domain.NewPortfolio(domain.NewPortfolioID("test-portfolio"), "test-user")

	// 異なる戦略の投資を追加
	investments := []struct {
		amount   float64
		strategy domain.InvestmentStrategy
	}{
		{1000000, domain.Conservative},
		{2000000, domain.Moderate},
		{3000000, domain.Aggressive},
	}

	for i, inv := range investments {
		money, _ := domain.NewMoney(inv.amount, "JPY")
		investment, _ := domain.NewInvestment(
			domain.NewInvestmentID(fmt.Sprintf("test-investment-%d", i)),
			money,
			domain.Stock,
			inv.strategy,
		)
		portfolio.AddInvestment(investment)
	}

	portfolioRepo.Save(ctx, portfolio)

	tests := []struct {
		name        string
		portfolioID string
		expectError bool
	}{
		{
			name:        "get analysis for existing portfolio",
			portfolioID: "test-portfolio",
			expectError: false,
		},
		{
			name:        "get analysis for non-existent portfolio",
			portfolioID: "non-existent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis, err := useCase.GetPortfolioAnalysis(ctx, tt.portfolioID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if analysis == nil {
					t.Error("Expected analysis, got nil")
				}
				if analysis.Portfolio.ID().Value != tt.portfolioID {
					t.Errorf("Expected portfolio ID %s, got %s", tt.portfolioID, analysis.Portfolio.ID().Value)
				}

				// ポートフォリオの合計金額を検証
				expectedTotal := 6000000.0 // 1M + 2M + 3M
				if analysis.TotalAmount.Amount != expectedTotal {
					t.Errorf("Expected total amount %f, got %f", expectedTotal, analysis.TotalAmount.Amount)
				}

				// リスクスコアを検証
				if analysis.RiskScore <= 0 || analysis.RiskScore > 1 {
					t.Errorf("Risk score %f is out of valid range (0,1]", analysis.RiskScore)
				}

				// 戦略配分を検証
				if len(analysis.StrategyAllocation) != 3 {
					t.Error("Expected allocations for all three strategies")
				}
			}
		})
	}
}

func TestPortfolioUseCase_RebalancePortfolio(t *testing.T) {
	ctx := context.Background()
	portfolioRepo := newMockPortfolioRepository()
	txManager := &mockTransactionManager{}
	eventPublisher := &mockEventPublisher{}
	strategyService := service.NewInvestmentStrategyService()

	useCase := NewPortfolioUseCase(
		portfolioRepo,
		txManager,
		eventPublisher,
		strategyService,
	)

	// テスト用のポートフォリオを作成
	portfolio := domain.NewPortfolio(domain.NewPortfolioID("test-portfolio"), "test-user")

	// 初期投資を追加（アグレッシブ戦略の金額を減らす）
	money, _ := domain.NewMoney(2000000, "JPY")
	investment, _ := domain.NewInvestment(
		domain.NewInvestmentID("test-investment"),
		money,
		domain.Stock,
		domain.Moderate, // より穏当な戦略に変更
	)
	portfolio.AddInvestment(investment)
	portfolioRepo.Save(ctx, portfolio)

	// リバランス用の新しい金額（より小さい金額に）
	newMoney, _ := domain.NewMoney(1500000, "JPY")
	changes := map[domain.InvestmentID]domain.Money{
		investment.ID(): newMoney,
	}

	tests := []struct {
		name        string
		portfolioID string
		changes     map[domain.InvestmentID]domain.Money
		expectError bool
	}{
		{
			name:        "valid rebalancing",
			portfolioID: "test-portfolio",
			changes:     changes,
			expectError: false,
		},
		{
			name:        "rebalance non-existent portfolio",
			portfolioID: "non-existent",
			changes:     changes,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := useCase.RebalancePortfolio(ctx, tt.portfolioID, tt.changes)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// リバランス後のポートフォリオを検証
				updatedPortfolio, _ := portfolioRepo.FindByID(ctx, domain.NewPortfolioID(tt.portfolioID))
				for invID, newAmount := range tt.changes {
					investment, err := updatedPortfolio.GetInvestment(invID)
					if err != nil {
						t.Errorf("Failed to get investment after rebalancing: %v", err)
						continue
					}

					if investment.Amount().Amount != newAmount.Amount {
						t.Errorf("Expected investment amount %f, got %f", newAmount.Amount, investment.Amount().Amount)
					}
				}
			}
		})
	}
}

type mockPortfolioRepository struct {
	portfolios map[domain.PortfolioID]*domain.Portfolio
}

func newMockPortfolioRepository() *mockPortfolioRepository {
	return &mockPortfolioRepository{
		portfolios: make(map[domain.PortfolioID]*domain.Portfolio),
	}
}

func (m *mockPortfolioRepository) Create(ctx context.Context, portfolio *domain.Portfolio) error {
	m.portfolios[portfolio.ID()] = portfolio
	return nil
}

func (m *mockPortfolioRepository) Save(ctx context.Context, portfolio *domain.Portfolio) error {
	m.portfolios[portfolio.ID()] = portfolio
	return nil
}

func (m *mockPortfolioRepository) FindByID(ctx context.Context, id domain.PortfolioID) (*domain.Portfolio, error) {
	if p, exists := m.portfolios[id]; exists {
		return p, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockPortfolioRepository) FindByUserID(ctx context.Context, userID string) (*domain.Portfolio, error) {
	for _, p := range m.portfolios {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockPortfolioRepository) FindByInvestmentID(ctx context.Context, investmentID domain.InvestmentID) (*domain.Portfolio, error) {
	return nil, domain.ErrNotFound
}

func (m *mockPortfolioRepository) Delete(ctx context.Context, id domain.PortfolioID) error {
	delete(m.portfolios, id)
	return nil
}

func (m *mockPortfolioRepository) Update(ctx context.Context, portfolio *domain.Portfolio) error {
	m.portfolios[portfolio.ID()] = portfolio
	return nil
}
