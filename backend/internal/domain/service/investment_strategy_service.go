package service

import (
	"errors"
	"moneyget/internal/domain"
)

type InvestmentStrategyService struct {
	maxAggressiveRatio  float64
	maxInvestmentAmount domain.Money
}

func NewInvestmentStrategyService() *InvestmentStrategyService {
	maxAmount, _ := domain.NewMoney(10000000, "JPY") // 1000万円
	return &InvestmentStrategyService{
		maxAggressiveRatio:  0.5, // 50%
		maxInvestmentAmount: maxAmount,
	}
}

func (s *InvestmentStrategyService) ValidateInvestmentStrategy(
	investment *domain.Investment,
	portfolio *domain.Portfolio,
) error {
	// 投資額の上限チェック
	totalAmount := portfolio.CalculateTotalAmount()
	newAmount, err := totalAmount.Add(investment.Amount())
	if err != nil {
		return err
	}
	if newAmount.Amount > s.maxInvestmentAmount.Amount {
		return domain.ErrPortfolioLimitExceeded
	}

	// アグレッシブ投資の比率チェック
	if investment.Strategy() == domain.Aggressive {
		currentAggressiveAmount := portfolio.CalculateStrategyAmount(domain.Aggressive)
		newAggressiveAmount, err := currentAggressiveAmount.Add(investment.Amount())
		if err != nil {
			return err
		}
		aggressiveRatio := newAggressiveAmount.Amount / newAmount.Amount
		if aggressiveRatio > s.maxAggressiveRatio {
			return domain.ErrAggressiveInvestmentLimitExceeded
		}
	}
	return nil
}

func (s *InvestmentStrategyService) CalculateRiskScore(portfolio *domain.Portfolio) (float64, error) {
	if portfolio == nil {
		return 0, errors.New("portfolio cannot be nil")
	}

	riskScores := map[domain.InvestmentStrategy]float64{
		domain.Conservative: 0.2,
		domain.Moderate:     0.5,
		domain.Aggressive:   1.0,
	}

	totalAmount := portfolio.CalculateTotalAmount()
	if totalAmount.Amount == 0 {
		return 0, nil
	}

	var weightedRiskScore float64
	investments := portfolio.GetInvestments()
	for _, investment := range investments {
		amount := investment.Amount()
		ratio := amount.Amount / totalAmount.Amount
		weightedRiskScore += ratio * riskScores[investment.Strategy()]
	}

	return weightedRiskScore, nil
}

func (s *InvestmentStrategyService) SuggestRebalancing(portfolio *domain.Portfolio) ([]RebalancingSuggestion, error) {
	if portfolio == nil {
		return nil, errors.New("portfolio cannot be nil")
	}

	var suggestions []RebalancingSuggestion
	allocation := make(map[domain.InvestmentStrategy]float64)
	totalAmount := portfolio.CalculateTotalAmount()

	if totalAmount.Amount == 0 {
		return suggestions, nil
	}

	// 各戦略の配分を計算
	for _, investment := range portfolio.GetInvestments() {
		strategy := investment.Strategy()
		amount := investment.Amount()
		allocation[strategy] = allocation[strategy] + (amount.Amount / totalAmount.Amount)
	}

	// アグレッシブ投資の比率チェック
	if aggressiveRatio := allocation[domain.Aggressive]; aggressiveRatio > s.maxAggressiveRatio {
		suggestions = append(suggestions, RebalancingSuggestion{
			Action:   "REDUCE",
			Strategy: domain.Aggressive,
			Reason:   "Aggressive allocation exceeds recommended maximum",
		})
	}

	// コンサバティブ投資の最小比率チェック
	if conservativeRatio := allocation[domain.Conservative]; conservativeRatio < 0.2 {
		suggestions = append(suggestions, RebalancingSuggestion{
			Action:   "INCREASE",
			Strategy: domain.Conservative,
			Reason:   "Conservative allocation below recommended minimum",
		})
	}

	return suggestions, nil
}

type RebalancingSuggestion struct {
	Action   string
	Strategy domain.InvestmentStrategy
	Reason   string
}
