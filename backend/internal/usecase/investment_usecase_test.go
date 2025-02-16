package usecase

import (
	"context"
	"moneyget/internal/domain"
	"moneyget/internal/domain/service"
	"testing"
)

// モックの定義
type mockInvestmentRepository struct {
	investments map[domain.InvestmentID]*domain.Investment
}

func newMockInvestmentRepository() *mockInvestmentRepository {
	return &mockInvestmentRepository{
		investments: make(map[domain.InvestmentID]*domain.Investment),
	}
}

func (m *mockInvestmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	m.investments[investment.ID()] = investment
	return nil
}

func (m *mockInvestmentRepository) Save(ctx context.Context, investment *domain.Investment) error {
	m.investments[investment.ID()] = investment
	return nil
}

func (m *mockInvestmentRepository) FindByID(ctx context.Context, id domain.InvestmentID) (*domain.Investment, error) {
	if inv, exists := m.investments[id]; exists {
		return inv, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockInvestmentRepository) FindAllByPortfolioID(ctx context.Context, portfolioID domain.PortfolioID) ([]*domain.Investment, error) {
	var result []*domain.Investment
	for _, inv := range m.investments {
		result = append(result, inv)
	}
	return result, nil
}

func (m *mockInvestmentRepository) Delete(ctx context.Context, id domain.InvestmentID) error {
	delete(m.investments, id)
	return nil
}

func (m *mockInvestmentRepository) FindAll(ctx context.Context) ([]*domain.Investment, error) {
	var result []*domain.Investment
	for _, inv := range m.investments {
		result = append(result, inv)
	}
	return result, nil
}

type mockTransactionManager struct{}

func (m *mockTransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockEventPublisher struct{}

func (m *mockEventPublisher) Publish(event domain.DomainEvent) error {
	return nil
}

func (m *mockEventPublisher) Subscribe(handler func(domain.DomainEvent)) error {
	return nil
}

func TestInvestmentUseCase_CreateInvestment(t *testing.T) {
	ctx := context.Background()
	investmentRepo := newMockInvestmentRepository()
	portfolioRepo := newPortfolioRepositoryForTest()
	txManager := &mockTransactionManager{}
	eventPublisher := &mockEventPublisher{}
	strategyService := service.NewInvestmentStrategyService()

	useCase := NewInvestmentUseCase(
		investmentRepo,
		portfolioRepo,
		txManager,
		eventPublisher,
		strategyService,
	)

	// ユーザーのポートフォリオを作成
	portfolio := domain.NewPortfolio(domain.NewPortfolioID("test-portfolio"), "test-user")
	portfolioRepo.Save(ctx, portfolio)

	tests := []struct {
		name        string
		userID      string
		amount      float64
		currency    string
		invType     string
		strategy    string
		expectError bool
	}{
		{
			name:        "valid investment creation",
			userID:      "test-user",
			amount:      1000000,
			currency:    "JPY",
			invType:     string(domain.Stock),
			strategy:    string(domain.Conservative),
			expectError: false,
		},
		{
			name:        "invalid investment type",
			userID:      "test-user",
			amount:      1000000,
			currency:    "JPY",
			invType:     "INVALID",
			strategy:    string(domain.Conservative),
			expectError: true,
		},
		{
			name:        "invalid amount",
			userID:      "test-user",
			amount:      -1000,
			currency:    "JPY",
			invType:     string(domain.Stock),
			strategy:    string(domain.Conservative),
			expectError: true,
		},
		{
			name:        "invalid currency",
			userID:      "test-user",
			amount:      1000000,
			currency:    "INVALID",
			invType:     string(domain.Stock),
			strategy:    string(domain.Conservative),
			expectError: true,
		},
		{
			name:        "portfolio not found",
			userID:      "non-existent-user",
			amount:      1000000,
			currency:    "JPY",
			invType:     string(domain.Stock),
			strategy:    string(domain.Conservative),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := useCase.CreateInvestment(
				ctx,
				tt.userID,
				tt.amount,
				tt.currency,
				tt.invType,
				tt.strategy,
			)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestInvestmentUseCase_GetInvestment(t *testing.T) {
	ctx := context.Background()
	investmentRepo := newMockInvestmentRepository()
	portfolioRepo := newPortfolioRepositoryForTest()
	txManager := &mockTransactionManager{}
	eventPublisher := &mockEventPublisher{}
	strategyService := service.NewInvestmentStrategyService()

	useCase := NewInvestmentUseCase(
		investmentRepo,
		portfolioRepo,
		txManager,
		eventPublisher,
		strategyService,
	)

	// テスト用の投資を作成
	money, _ := domain.NewMoney(1000000, "JPY")
	investment, _ := domain.NewInvestment(
		domain.NewInvestmentID("test-investment"),
		money,
		domain.Stock,
		domain.Conservative,
	)
	investmentRepo.Save(ctx, investment)

	tests := []struct {
		name        string
		id          string
		expectError bool
	}{
		{
			name:        "get existing investment",
			id:          "test-investment",
			expectError: false,
		},
		{
			name:        "get non-existent investment",
			id:          "non-existent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, err := useCase.GetInvestment(ctx, tt.id)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if inv == nil {
					t.Error("Expected investment, got nil")
				}
				if inv.ID().Value != tt.id {
					t.Errorf("Expected investment ID %s, got %s", tt.id, inv.ID().Value)
				}
			}
		})
	}
}

// テスト用のPortfolioRepositoryを取得する関数
func newPortfolioRepositoryForTest() domain.PortfolioRepository {
	return &portfolioRepoFromTest{
		portfolios: make(map[domain.PortfolioID]*domain.Portfolio),
	}
}

type portfolioRepoFromTest struct {
	portfolios map[domain.PortfolioID]*domain.Portfolio
}

func (m *portfolioRepoFromTest) Create(ctx context.Context, portfolio *domain.Portfolio) error {
	m.portfolios[portfolio.ID()] = portfolio
	return nil
}

func (m *portfolioRepoFromTest) Save(ctx context.Context, portfolio *domain.Portfolio) error {
	m.portfolios[portfolio.ID()] = portfolio
	return nil
}

func (m *portfolioRepoFromTest) FindByID(ctx context.Context, id domain.PortfolioID) (*domain.Portfolio, error) {
	if p, exists := m.portfolios[id]; exists {
		return p, nil
	}
	return nil, domain.ErrNotFound
}

func (m *portfolioRepoFromTest) FindByUserID(ctx context.Context, userID string) (*domain.Portfolio, error) {
	for _, p := range m.portfolios {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *portfolioRepoFromTest) FindByInvestmentID(ctx context.Context, investmentID domain.InvestmentID) (*domain.Portfolio, error) {
	return nil, domain.ErrNotFound
}

func (m *portfolioRepoFromTest) Delete(ctx context.Context, id domain.PortfolioID) error {
	delete(m.portfolios, id)
	return nil
}

func (m *portfolioRepoFromTest) Update(ctx context.Context, portfolio *domain.Portfolio) error {
	m.portfolios[portfolio.ID()] = portfolio
	return nil
}
