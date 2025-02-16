package usecase

import (
	"context"
	"moneyget/internal/domain"
	"moneyget/internal/domain/service"
	"moneyget/internal/utils"
)

type InvestmentUseCase struct {
	investmentRepo  domain.InvestmentRepository
	portfolioRepo   domain.PortfolioRepository
	txManager       domain.TransactionManager
	eventPublisher  domain.DomainEventPublisher
	strategyService *service.InvestmentStrategyService
}

func NewInvestmentUseCase(
	investmentRepo domain.InvestmentRepository,
	portfolioRepo domain.PortfolioRepository,
	txManager domain.TransactionManager,
	eventPublisher domain.DomainEventPublisher,
	strategyService *service.InvestmentStrategyService,
) *InvestmentUseCase {
	return &InvestmentUseCase{
		investmentRepo:  investmentRepo,
		portfolioRepo:   portfolioRepo,
		txManager:       txManager,
		eventPublisher:  eventPublisher,
		strategyService: strategyService,
	}
}

// CreateInvestmentメソッドのシグネチャを変更
func (u *InvestmentUseCase) CreateInvestment(
	ctx context.Context,
	userID string,
	amount float64,
	currency string,
	investmentType string,
	strategy string,
) error {
	return u.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		// ユーザーのポートフォリオを取得
		portfolio, err := u.portfolioRepo.FindByUserID(ctx, userID)
		if err != nil {
			return err
		}

		money, err := domain.NewMoney(amount, currency)
		if err != nil {
			return err
		}

		investment, err := domain.NewInvestment(
			domain.NewInvestmentID(utils.GenerateUUID()),
			money,
			domain.InvestmentType(investmentType),
			domain.InvestmentStrategy(strategy),
		)
		if err != nil {
			return err
		}

		// 投資戦略の検証
		if err := u.strategyService.ValidateInvestmentStrategy(investment, portfolio); err != nil {
			return err
		}

		if err := portfolio.AddInvestment(investment); err != nil {
			return err
		}

		if err := u.investmentRepo.Save(ctx, investment); err != nil {
			return err
		}

		if err := u.portfolioRepo.Save(ctx, portfolio); err != nil {
			return err
		}

		event := domain.NewInvestmentCreatedEvent(investment.ID(), money)
		return u.eventPublisher.Publish(event)
	})
}

func (u *InvestmentUseCase) GetInvestment(
	ctx context.Context,
	id string,
) (*domain.Investment, error) {
	return u.investmentRepo.FindByID(ctx, domain.NewInvestmentID(id))
}

func (u *InvestmentUseCase) GetPortfolioInvestments(
	ctx context.Context,
	portfolioID string,
) ([]*domain.Investment, error) {
	return u.investmentRepo.FindAllByPortfolioID(ctx, domain.NewPortfolioID(portfolioID))
}

type InvestmentWithRisk struct {
	Investment *domain.Investment
	RiskScore  float64
}

func (u *InvestmentUseCase) GetInvestmentWithRiskAnalysis(
	ctx context.Context,
	id string,
) (*InvestmentWithRisk, error) {
	investment, err := u.investmentRepo.FindByID(ctx, domain.NewInvestmentID(id))
	if err != nil {
		return nil, err
	}

	portfolio, err := u.findPortfolioByInvestmentID(ctx, id)
	if err != nil {
		return nil, err
	}

	riskScore, err := u.strategyService.CalculateRiskScore(portfolio)
	if err != nil {
		return nil, err
	}

	return &InvestmentWithRisk{
		Investment: investment,
		RiskScore:  riskScore,
	}, nil
}

func (u *InvestmentUseCase) GetInvestmentRebalancingSuggestions(
	ctx context.Context,
	portfolioID string,
) ([]service.RebalancingSuggestion, error) {
	portfolio, err := u.portfolioRepo.FindByID(ctx, domain.NewPortfolioID(portfolioID))
	if err != nil {
		return nil, err
	}

	return u.strategyService.SuggestRebalancing(portfolio)
}

func (u *InvestmentUseCase) findPortfolioByInvestmentID(
	ctx context.Context,
	investmentID string,
) (*domain.Portfolio, error) {
	return u.portfolioRepo.FindByInvestmentID(ctx, domain.NewInvestmentID(investmentID))
}
