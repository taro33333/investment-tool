package usecase

import (
	"context"
	"moneyget/internal/domain"
	"moneyget/internal/domain/service"
	"moneyget/internal/utils"
)

type PortfolioUseCase struct {
	portfolioRepo   domain.PortfolioRepository
	txManager       domain.TransactionManager
	eventPublisher  domain.DomainEventPublisher
	strategyService *service.InvestmentStrategyService
}

func NewPortfolioUseCase(
	portfolioRepo domain.PortfolioRepository,
	txManager domain.TransactionManager,
	eventPublisher domain.DomainEventPublisher,
	strategyService *service.InvestmentStrategyService,
) *PortfolioUseCase {
	return &PortfolioUseCase{
		portfolioRepo:   portfolioRepo,
		txManager:       txManager,
		eventPublisher:  eventPublisher,
		strategyService: strategyService,
	}
}

type PortfolioAnalysis struct {
	Portfolio          *domain.Portfolio
	TotalAmount        domain.Money
	RiskScore          float64
	StrategyAllocation map[domain.InvestmentStrategy]float64
	Suggestions        []service.RebalancingSuggestion
}

func (u *PortfolioUseCase) GetPortfolioAnalysis(ctx context.Context, id string) (*PortfolioAnalysis, error) {
	portfolio, err := u.portfolioRepo.FindByID(ctx, domain.NewPortfolioID(id))
	if err != nil {
		return nil, err
	}

	riskScore, err := u.strategyService.CalculateRiskScore(portfolio)
	if err != nil {
		return nil, err
	}

	suggestions, err := u.strategyService.SuggestRebalancing(portfolio)
	if err != nil {
		return nil, err
	}

	allocation := u.calculateStrategyAllocation(portfolio)

	return &PortfolioAnalysis{
		Portfolio:          portfolio,
		TotalAmount:        portfolio.CalculateTotalAmount(),
		RiskScore:          riskScore,
		StrategyAllocation: allocation,
		Suggestions:        suggestions,
	}, nil
}

func (u *PortfolioUseCase) RebalancePortfolio(ctx context.Context, id string, changes map[domain.InvestmentID]domain.Money) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		portfolio, err := u.portfolioRepo.FindByID(ctx, domain.NewPortfolioID(id))
		if err != nil {
			return err
		}

		for investmentID, newAmount := range changes {
			investment, err := portfolio.GetInvestment(investmentID)
			if err != nil {
				return err
			}

			if err := investment.UpdateAmount(newAmount); err != nil {
				return err
			}
		}

		// 再配分後の検証
		if err := portfolio.ValidateRiskDistribution(); err != nil {
			return err
		}

		if err := u.portfolioRepo.Save(ctx, portfolio); err != nil {
			return err
		}

		event := domain.NewPortfolioUpdatedEvent(portfolio.ID(), portfolio.CalculateTotalAmount())
		return u.eventPublisher.Publish(event)
	})
}

func (u *PortfolioUseCase) calculateStrategyAllocation(portfolio *domain.Portfolio) map[domain.InvestmentStrategy]float64 {
	allocation := make(map[domain.InvestmentStrategy]float64)
	totalAmount := portfolio.CalculateTotalAmount()

	if totalAmount.Amount == 0 {
		return allocation
	}

	for _, investment := range portfolio.GetInvestments() {
		strategy := investment.Strategy()
		currentRatio := allocation[strategy]
		allocation[strategy] = currentRatio + (investment.Amount().Amount / totalAmount.Amount)
	}

	return allocation
}

func (u *PortfolioUseCase) CreatePortfolio(ctx context.Context, userID string) (*domain.Portfolio, error) {
	portfolio := domain.NewPortfolio(domain.NewPortfolioID(utils.GenerateUUID()), userID)

	err := u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := u.portfolioRepo.Save(ctx, portfolio); err != nil {
			return err
		}

		totalAmount, _ := domain.NewMoney(0, "JPY")
		event := domain.NewPortfolioUpdatedEvent(portfolio.ID(), totalAmount)
		return u.eventPublisher.Publish(event)
	})

	if err != nil {
		return nil, err
	}

	return portfolio, nil
}

func (u *PortfolioUseCase) GetPortfolio(ctx context.Context, id string) (*domain.Portfolio, error) {
	return u.portfolioRepo.FindByID(ctx, domain.NewPortfolioID(id))
}

func (u *PortfolioUseCase) GetUserPortfolio(ctx context.Context, userID string) (*domain.Portfolio, error) {
	return u.portfolioRepo.FindByUserID(ctx, userID)
}

func (u *PortfolioUseCase) ValidatePortfolio(ctx context.Context, portfolioID string) error {
	portfolio, err := u.portfolioRepo.FindByID(ctx, domain.NewPortfolioID(portfolioID))
	if err != nil {
		return err
	}

	return portfolio.ValidateRiskDistribution()
}

func generateUUID() string {
	// UUIDの生成ロジックを実装
	// 実際のプロジェクトではgithub.com/google/uuidなどのライブラリを使用することを推奨
	return "uuid-implementation-needed"
}
