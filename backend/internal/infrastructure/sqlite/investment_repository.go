package sqlite

import (
	"context"
	"database/sql"
	"moneyget/internal/domain"
)

type investmentRepository struct {
	db *sql.DB
}

func NewInvestmentRepository(db *sql.DB) domain.InvestmentRepository {
	return &investmentRepository{db: db}
}

func (r *investmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	query := `
		INSERT INTO investments (id, amount, currency, type, strategy, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	amount := investment.Amount()
	_, err := r.db.ExecContext(ctx, query,
		investment.ID().Value,
		amount.Amount,
		amount.Currency,
		string(investment.Type()),
		string(investment.Strategy()),
		investment.CreatedAt,
		investment.UpdatedAt,
	)
	return err
}

func (r *investmentRepository) Save(ctx context.Context, investment *domain.Investment) error {
	query := `
		INSERT INTO investments (id, amount, currency, type, strategy, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			amount = excluded.amount,
			currency = excluded.currency,
			type = excluded.type,
			strategy = excluded.strategy,
			updated_at = excluded.updated_at
	`

	amount := investment.Amount()
	_, err := r.db.ExecContext(ctx, query,
		investment.ID().Value,
		amount.Amount,
		amount.Currency,
		string(investment.Type()),
		string(investment.Strategy()),
		investment.CreatedAt,
		investment.UpdatedAt,
	)
	return err
}

func (r *investmentRepository) FindByID(ctx context.Context, id domain.InvestmentID) (*domain.Investment, error) {
	query := `
		SELECT amount, currency, type, strategy, created_at, updated_at
		FROM investments
		WHERE id = ?
	`

	var amount float64
	var currency string
	var investmentType string
	var strategy string
	var createdAt string
	var updatedAt string

	err := r.db.QueryRowContext(ctx, query, id.Value).Scan(
		&amount,
		&currency,
		&investmentType,
		&strategy,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	money, err := domain.NewMoney(amount, currency)
	if err != nil {
		return nil, err
	}

	return domain.NewInvestment(
		id,
		money,
		domain.InvestmentType(investmentType),
		domain.InvestmentStrategy(strategy),
	)
}

func (r *investmentRepository) FindAllByPortfolioID(ctx context.Context, portfolioID domain.PortfolioID) ([]*domain.Investment, error) {
	query := `
		SELECT i.id, i.amount, i.currency, i.type, i.strategy, i.created_at, i.updated_at
		FROM investments i
		JOIN portfolio_investments pi ON i.id = pi.investment_id
		WHERE pi.portfolio_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, portfolioID.Value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []*domain.Investment
	for rows.Next() {
		var id string
		var amount float64
		var currency string
		var investmentType string
		var strategy string
		var createdAt string
		var updatedAt string

		err := rows.Scan(&id, &amount, &currency, &investmentType, &strategy, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		money, err := domain.NewMoney(amount, currency)
		if err != nil {
			return nil, err
		}

		investment, err := domain.NewInvestment(
			domain.NewInvestmentID(id),
			money,
			domain.InvestmentType(investmentType),
			domain.InvestmentStrategy(strategy),
		)
		if err != nil {
			return nil, err
		}

		investments = append(investments, investment)
	}

	return investments, nil
}

func (r *investmentRepository) FindAll(ctx context.Context) ([]*domain.Investment, error) {
	query := `
		SELECT id, amount, currency, type, strategy, created_at, updated_at
		FROM investments
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var investments []*domain.Investment
	for rows.Next() {
		var id string
		var amount float64
		var currency string
		var investmentType string
		var strategy string
		var createdAt string
		var updatedAt string

		err := rows.Scan(&id, &amount, &currency, &investmentType, &strategy, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		money, err := domain.NewMoney(amount, currency)
		if err != nil {
			return nil, err
		}

		investment, err := domain.NewInvestment(
			domain.NewInvestmentID(id),
			money,
			domain.InvestmentType(investmentType),
			domain.InvestmentStrategy(strategy),
		)
		if err != nil {
			return nil, err
		}

		investments = append(investments, investment)
	}

	return investments, nil
}

func (r *investmentRepository) Delete(ctx context.Context, id domain.InvestmentID) error {
	query := "DELETE FROM investments WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id.Value)
	return err
}
