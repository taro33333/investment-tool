package domain

import (
	"errors"
	"time"
)

type PortfolioID struct {
	Value string // エクスポート
}

func NewPortfolioID(id string) PortfolioID {
	return PortfolioID{Value: id}
}

type Portfolio struct {
	id          PortfolioID
	UserID      string                       // エクスポート
	Investments map[InvestmentID]*Investment // エクスポート
	CreatedAt   time.Time                    // エクスポート
	UpdatedAt   time.Time                    // エクスポート
}

func NewPortfolio(id PortfolioID, userID string) *Portfolio {
	now := time.Now()
	return &Portfolio{
		id:          id,
		UserID:      userID,
		Investments: make(map[InvestmentID]*Investment),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Portfolio) ID() PortfolioID {
	return p.id
}

func (p *Portfolio) AddInvestment(investment *Investment) error {
	if investment == nil {
		return errors.New("investment cannot be nil")
	}

	if _, exists := p.Investments[investment.ID()]; exists {
		return ErrDuplicateInvestment
	}

	totalAmount := p.CalculateTotalAmount()
	newAmount := totalAmount.Amount + investment.Amount().Amount
	if newAmount > 10000000 { // 1000万円の投資上限を設定
		return ErrPortfolioLimitExceeded
	}

	p.Investments[investment.ID()] = investment
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Portfolio) RemoveInvestment(investmentID InvestmentID) error {
	if _, exists := p.Investments[investmentID]; !exists {
		return ErrInvestmentNotFound
	}

	delete(p.Investments, investmentID)
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Portfolio) GetInvestment(investmentID InvestmentID) (*Investment, error) {
	investment, exists := p.Investments[investmentID]
	if !exists {
		return nil, ErrInvestmentNotFound
	}
	return investment, nil
}

func (p *Portfolio) GetInvestments() []*Investment {
	investments := make([]*Investment, 0, len(p.Investments))
	for _, inv := range p.Investments {
		investments = append(investments, inv)
	}
	return investments
}

func (p *Portfolio) CalculateTotalAmount() Money {
	var total float64
	currency := "JPY" // デフォルト通貨を設定
	for _, inv := range p.Investments {
		total += inv.Amount().Amount
		currency = inv.Amount().Currency // 既存の投資がある場合はその通貨を使用
	}
	money, _ := NewMoney(total, currency)
	return money
}

func (p *Portfolio) CalculateStrategyAmount(strategy InvestmentStrategy) Money {
	var total float64
	var currency string
	for _, inv := range p.Investments {
		if inv.Strategy() == strategy {
			total += inv.Amount().Amount
			currency = inv.Amount().Currency
		}
	}
	money, _ := NewMoney(total, currency)
	return money
}

func (p *Portfolio) ValidateRiskDistribution() error {
	var aggressiveTotal float64
	totalAmount := p.CalculateTotalAmount()

	if totalAmount.Amount == 0 {
		return nil
	}

	for _, inv := range p.Investments {
		if inv.Strategy() == Aggressive {
			aggressiveTotal += inv.Amount().Amount
		}
	}

	aggressiveRatio := aggressiveTotal / totalAmount.Amount
	if aggressiveRatio > 0.5 { // ポートフォリオの50%以上をアグレッシブ投資にはできない
		return ErrAggressiveInvestmentLimitExceeded
	}

	return nil
}
