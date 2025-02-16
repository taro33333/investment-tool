package domain

import (
	"testing"
)

func TestNewPortfolio(t *testing.T) {
	id := NewPortfolioID("test-portfolio")
	userID := "test-user"
	portfolio := NewPortfolio(id, userID)

	if portfolio.ID() != id {
		t.Errorf("Expected portfolio ID %v, got %v", id, portfolio.ID())
	}

	if portfolio.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, portfolio.UserID)
	}

	if len(portfolio.GetInvestments()) != 0 {
		t.Error("New portfolio should have no investments")
	}
}

func TestPortfolio_AddInvestment(t *testing.T) {
	portfolio := NewPortfolio(NewPortfolioID("test-portfolio"), "test-user")

	tests := []struct {
		name        string
		amount      float64
		strategy    InvestmentStrategy
		expectError bool
	}{
		{
			name:        "add valid investment",
			amount:      1000000,
			strategy:    Conservative,
			expectError: false,
		},
		{
			name:        "exceed portfolio limit",
			amount:      11000000,
			strategy:    Conservative,
			expectError: true,
		},
		{
			name:        "add duplicate investment",
			amount:      1000000,
			strategy:    Conservative,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, _ := NewMoney(tt.amount, "JPY")
			investment, _ := NewInvestment(
				NewInvestmentID("test-investment"),
				money,
				Stock,
				tt.strategy,
			)

			err := portfolio.AddInvestment(investment)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify investment was added
				inv, err := portfolio.GetInvestment(investment.ID())
				if err != nil {
					t.Errorf("Failed to get added investment: %v", err)
				}
				if inv != investment {
					t.Error("Retrieved investment does not match added investment")
				}
			}
		})
	}
}

func TestPortfolio_RemoveInvestment(t *testing.T) {
	portfolio := NewPortfolio(NewPortfolioID("test-portfolio"), "test-user")

	// Add an investment
	money, _ := NewMoney(1000000, "JPY")
	investment, _ := NewInvestment(
		NewInvestmentID("test-investment"),
		money,
		Stock,
		Conservative,
	)
	_ = portfolio.AddInvestment(investment)

	// Test removal
	err := portfolio.RemoveInvestment(investment.ID())
	if err != nil {
		t.Errorf("Unexpected error removing investment: %v", err)
	}

	// Verify investment was removed
	_, err = portfolio.GetInvestment(investment.ID())
	if err != ErrInvestmentNotFound {
		t.Errorf("Expected ErrInvestmentNotFound, got %v", err)
	}

	// Test removing non-existent investment
	err = portfolio.RemoveInvestment(NewInvestmentID("non-existent"))
	if err != ErrInvestmentNotFound {
		t.Errorf("Expected ErrInvestmentNotFound for non-existent investment, got %v", err)
	}
}

func TestPortfolio_CalculateTotalAmount(t *testing.T) {
	portfolio := NewPortfolio(NewPortfolioID("test-portfolio"), "test-user")

	// Add multiple investments
	investments := []struct {
		amount   float64
		strategy InvestmentStrategy
	}{
		{1000000, Conservative},
		{2000000, Moderate},
		{3000000, Aggressive},
	}

	expectedTotal := 0.0
	for _, inv := range investments {
		money, _ := NewMoney(inv.amount, "JPY")
		investment, _ := NewInvestment(
			NewInvestmentID("test-investment-"+string(inv.strategy)),
			money,
			Stock,
			inv.strategy,
		)
		_ = portfolio.AddInvestment(investment)
		expectedTotal += inv.amount
	}

	totalAmount := portfolio.CalculateTotalAmount()
	if totalAmount.Amount != expectedTotal {
		t.Errorf("Expected total amount %f, got %f", expectedTotal, totalAmount.Amount)
	}
	if totalAmount.Currency != "JPY" {
		t.Errorf("Expected currency JPY, got %s", totalAmount.Currency)
	}
}
