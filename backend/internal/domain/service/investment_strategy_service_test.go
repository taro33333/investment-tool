package service

import (
	"fmt"
	"moneyget/internal/domain"
	"testing"
)

func TestInvestmentStrategyService_ValidateInvestmentStrategy(t *testing.T) {
	service := NewInvestmentStrategyService()

	tests := []struct {
		name       string
		investment struct {
			amount   float64
			strategy domain.InvestmentStrategy
		}
		existingInvestments []struct {
			amount   float64
			strategy domain.InvestmentStrategy
		}
		expectError error
	}{
		{
			name: "valid conservative investment",
			investment: struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{
				amount:   1000000,
				strategy: domain.Conservative,
			},
			existingInvestments: nil,
			expectError:         nil,
		},
		{
			name: "exceed portfolio limit",
			investment: struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{
				amount:   6000000,
				strategy: domain.Conservative,
			},
			existingInvestments: []struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{
				{5000000, domain.Conservative},
			},
			expectError: domain.ErrPortfolioLimitExceeded,
		},
		{
			name: "exceed aggressive ratio",
			investment: struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{
				amount:   6000000,
				strategy: domain.Aggressive,
			},
			existingInvestments: []struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{
				{4000000, domain.Aggressive},
			},
			expectError: domain.ErrAggressiveInvestmentLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portfolio := domain.NewPortfolio(domain.NewPortfolioID("test-portfolio"), "test-user")

			// Add existing investments
			if tt.existingInvestments != nil {
				for i, inv := range tt.existingInvestments {
					money, err := domain.NewMoney(inv.amount, "JPY")
					if err != nil {
						t.Fatalf("Failed to create money: %v", err)
					}
					investment, err := domain.NewInvestment(
						domain.NewInvestmentID(fmt.Sprintf("existing-investment-%d", i)),
						money,
						domain.Stock,
						inv.strategy,
					)
					if err != nil {
						t.Fatalf("Failed to create investment: %v", err)
					}
					if err := portfolio.AddInvestment(investment); err != nil {
						t.Fatalf("Failed to add existing investment: %v", err)
					}
				}
			}

			// Create and validate new investment
			money, err := domain.NewMoney(tt.investment.amount, "JPY")
			if err != nil {
				t.Fatalf("Failed to create money: %v", err)
			}
			investment, err := domain.NewInvestment(
				domain.NewInvestmentID("test-investment"),
				money,
				domain.Stock,
				tt.investment.strategy,
			)
			if err != nil {
				t.Fatalf("Failed to create investment: %v", err)
			}

			err = service.ValidateInvestmentStrategy(investment, portfolio)

			if err != tt.expectError {
				t.Errorf("Expected error %v, got %v", tt.expectError, err)
			}
		})
	}
}

func TestInvestmentStrategyService_CalculateRiskScore(t *testing.T) {
	service := NewInvestmentStrategyService()
	portfolio := domain.NewPortfolio(domain.NewPortfolioID("test-portfolio"), "test-user")

	tests := []struct {
		name        string
		investments []struct {
			amount   float64
			strategy domain.InvestmentStrategy
		}
		expectedScore float64
	}{
		{
			name: "empty portfolio",
			investments: []struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{},
			expectedScore: 0,
		},
		{
			name: "conservative only",
			investments: []struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{
				{1000000, domain.Conservative},
			},
			expectedScore: 0.2,
		},
		{
			name: "mixed portfolio",
			investments: []struct {
				amount   float64
				strategy domain.InvestmentStrategy
			}{
				{1000000, domain.Conservative}, // 0.2 * 0.333 = 0.0666
				{1000000, domain.Moderate},     // 0.5 * 0.333 = 0.1665
				{1000000, domain.Aggressive},   // 1.0 * 0.333 = 0.333
			},
			expectedScore: 0.566, // 約0.566
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset portfolio for each test
			portfolio = domain.NewPortfolio(domain.NewPortfolioID("test-portfolio"), "test-user")

			// Add investments
			for i, inv := range tt.investments {
				money, _ := domain.NewMoney(inv.amount, "JPY")
				investment, _ := domain.NewInvestment(
					domain.NewInvestmentID(fmt.Sprintf("test-investment-%d", i)),
					money,
					domain.Stock,
					inv.strategy,
				)
				portfolio.AddInvestment(investment)
			}

			score, err := service.CalculateRiskScore(portfolio)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Allow for small floating point differences
			if score < tt.expectedScore-0.01 || score > tt.expectedScore+0.01 {
				t.Errorf("Expected risk score approximately %f, got %f", tt.expectedScore, score)
			}
		})
	}
}

func TestInvestmentStrategyService_SuggestRebalancing(t *testing.T) {
	service := NewInvestmentStrategyService()
	portfolio := domain.NewPortfolio(domain.NewPortfolioID("test-portfolio"), "test-user")

	// Add investments with unbalanced allocation
	investments := []struct {
		amount   float64
		strategy domain.InvestmentStrategy
	}{
		{7000000, domain.Aggressive},   // 70% aggressive - このケースは明確に上限を超えている
		{2000000, domain.Moderate},     // 20% moderate
		{1000000, domain.Conservative}, // 10% conservative
	}

	for i, inv := range investments {
		money, _ := domain.NewMoney(inv.amount, "JPY")
		investment, _ := domain.NewInvestment(
			domain.NewInvestmentID(fmt.Sprintf("test-investment-%d", i)),
			money,
			domain.Stock,
			inv.strategy,
		)
		if err := portfolio.AddInvestment(investment); err != nil {
			t.Fatalf("Failed to add investment: %v", err)
		}
	}

	suggestions, err := service.SuggestRebalancing(portfolio)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(suggestions) == 0 {
		t.Error("Expected rebalancing suggestions for unbalanced portfolio, got none")
	}

	// 各提案をチェック
	foundAggressiveReduction := false
	for _, suggestion := range suggestions {
		if suggestion.Action == "REDUCE" && suggestion.Strategy == domain.Aggressive {
			foundAggressiveReduction = true
			break
		}
	}

	if !foundAggressiveReduction {
		t.Error("Expected suggestion to reduce aggressive allocation")
	}
}
