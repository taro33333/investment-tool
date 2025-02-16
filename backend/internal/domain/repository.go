package domain

import "context"

type UserRepository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}

type InvestmentRepository interface {
	Create(ctx context.Context, investment *Investment) error
	Save(ctx context.Context, investment *Investment) error
	FindByID(ctx context.Context, id InvestmentID) (*Investment, error)
	FindAllByPortfolioID(ctx context.Context, portfolioID PortfolioID) ([]*Investment, error)
	Delete(ctx context.Context, id InvestmentID) error
	FindAll(ctx context.Context) ([]*Investment, error)
}

type PortfolioRepository interface {
	Create(ctx context.Context, portfolio *Portfolio) error
	Save(ctx context.Context, portfolio *Portfolio) error
	FindByID(ctx context.Context, id PortfolioID) (*Portfolio, error)
	FindByUserID(ctx context.Context, userID string) (*Portfolio, error)
	FindByInvestmentID(ctx context.Context, investmentID InvestmentID) (*Portfolio, error)
	Delete(ctx context.Context, id PortfolioID) error
	Update(ctx context.Context, portfolio *Portfolio) error
}

type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
