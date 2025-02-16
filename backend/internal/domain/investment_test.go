package domain

import (
	"testing"
	"time"
)

func TestNewInvestment(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		amount      float64
		currency    string
		typeVal     InvestmentType
		strategy    InvestmentStrategy
		expectError bool
	}{
		{
			name:        "valid investment creation",
			id:          "test-id",
			amount:      1000.0,
			currency:    "JPY",
			typeVal:     Stock,
			strategy:    Conservative,
			expectError: false,
		},
		{
			name:        "invalid investment type",
			id:          "test-id",
			amount:      1000.0,
			currency:    "JPY",
			typeVal:     "INVALID",
			strategy:    Conservative,
			expectError: true,
		},
		{
			name:        "invalid strategy",
			id:          "test-id",
			amount:      1000.0,
			currency:    "JPY",
			typeVal:     Stock,
			strategy:    "INVALID",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)
			if err != nil {
				t.Fatalf("Failed to create money: %v", err)
			}

			investment, err := NewInvestment(
				NewInvestmentID(tt.id),
				money,
				tt.typeVal,
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
				if investment.ID().Value != tt.id {
					t.Errorf("Expected ID %s, got %s", tt.id, investment.ID().Value)
				}
				if investment.Amount().Amount != tt.amount {
					t.Errorf("Expected amount %f, got %f", tt.amount, investment.Amount().Amount)
				}
				if investment.Type() != tt.typeVal {
					t.Errorf("Expected type %s, got %s", tt.typeVal, investment.Type())
				}
				if investment.Strategy() != tt.strategy {
					t.Errorf("Expected strategy %s, got %s", tt.strategy, investment.Strategy())
				}
			}
		})
	}
}

func TestInvestment_UpdateAmount(t *testing.T) {
	initialMoney, _ := NewMoney(1000.0, "JPY")
	investment, _ := NewInvestment(
		NewInvestmentID("test-id"),
		initialMoney,
		Stock,
		Conservative,
	)

	// Record initial update time
	initialUpdateTime := investment.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond)

	// Update amount
	newMoney, _ := NewMoney(2000.0, "JPY")
	err := investment.UpdateAmount(newMoney)

	if err != nil {
		t.Errorf("Unexpected error updating amount: %v", err)
	}

	if investment.Amount().Amount != 2000.0 {
		t.Errorf("Expected amount 2000.0, got %f", investment.Amount().Amount)
	}

	if investment.UpdatedAt.Equal(initialUpdateTime) {
		t.Error("UpdatedAt time was not changed")
	}
}
