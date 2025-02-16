package sqlite

import (
	"context"
	"database/sql"
	"moneyget/internal/domain"
)

type portfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) domain.PortfolioRepository {
	return &portfolioRepository{db: db}
}

func (r *portfolioRepository) Create(ctx context.Context, portfolio *domain.Portfolio) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO portfolios (id, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, query,
		portfolio.ID().Value,
		portfolio.UserID,
		portfolio.CreatedAt,
		portfolio.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// 投資との関連付けを登録
	for _, investment := range portfolio.GetInvestments() {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO portfolio_investments (portfolio_id, investment_id) VALUES (?, ?)",
			portfolio.ID().Value,
			investment.ID().Value,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *portfolioRepository) Save(ctx context.Context, portfolio *domain.Portfolio) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save portfolio
	query := `
		INSERT INTO portfolios (id, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			user_id = excluded.user_id,
			updated_at = excluded.updated_at
	`
	_, err = tx.ExecContext(ctx, query,
		portfolio.ID().Value,
		portfolio.UserID,
		portfolio.CreatedAt,
		portfolio.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Delete existing portfolio investments
	_, err = tx.ExecContext(ctx, "DELETE FROM portfolio_investments WHERE portfolio_id = ?", portfolio.ID().Value)
	if err != nil {
		return err
	}

	// Save portfolio investments
	for _, investment := range portfolio.GetInvestments() {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO portfolio_investments (portfolio_id, investment_id) VALUES (?, ?)",
			portfolio.ID().Value,
			investment.ID().Value,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *portfolioRepository) FindByID(ctx context.Context, id domain.PortfolioID) (*domain.Portfolio, error) {
	query := `
		SELECT user_id, created_at, updated_at
		FROM portfolios
		WHERE id = ?
	`

	var userID string
	var createdAt string
	var updatedAt string

	err := r.db.QueryRowContext(ctx, query, id.Value).Scan(
		&userID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	portfolio := domain.NewPortfolio(id, userID)

	// Load investments
	investments, err := r.loadPortfolioInvestments(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, investment := range investments {
		if err := portfolio.AddInvestment(investment); err != nil {
			return nil, err
		}
	}

	return portfolio, nil
}

func (r *portfolioRepository) FindByUserID(ctx context.Context, userID string) (*domain.Portfolio, error) {
	query := `
		SELECT id
		FROM portfolios
		WHERE user_id = ?
	`

	var id string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, domain.NewPortfolioID(id))
}

func (r *portfolioRepository) Delete(ctx context.Context, id domain.PortfolioID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete portfolio investments
	_, err = tx.ExecContext(ctx, "DELETE FROM portfolio_investments WHERE portfolio_id = ?", id.Value)
	if err != nil {
		return err
	}

	// Delete portfolio
	_, err = tx.ExecContext(ctx, "DELETE FROM portfolios WHERE id = ?", id.Value)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *portfolioRepository) Update(ctx context.Context, portfolio *domain.Portfolio) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update portfolio
	query := `
		UPDATE portfolios
		SET user_id = ?, updated_at = ?
		WHERE id = ?
	`
	_, err = tx.ExecContext(ctx, query,
		portfolio.UserID,
		portfolio.UpdatedAt,
		portfolio.ID().Value,
	)
	if err != nil {
		return err
	}

	// Update portfolio investments
	_, err = tx.ExecContext(ctx, "DELETE FROM portfolio_investments WHERE portfolio_id = ?", portfolio.ID().Value)
	if err != nil {
		return err
	}

	for _, investment := range portfolio.GetInvestments() {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO portfolio_investments (portfolio_id, investment_id) VALUES (?, ?)",
			portfolio.ID().Value,
			investment.ID().Value,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *portfolioRepository) loadPortfolioInvestments(ctx context.Context, portfolioID domain.PortfolioID) ([]*domain.Investment, error) {
	investmentRepo := NewInvestmentRepository(r.db)
	return investmentRepo.FindAllByPortfolioID(ctx, portfolioID)
}

func (r *portfolioRepository) FindByInvestmentID(ctx context.Context, investmentID domain.InvestmentID) (*domain.Portfolio, error) {
	query := `
		SELECT p.id, p.user_id, p.created_at, p.updated_at
		FROM portfolios p
		JOIN portfolio_investments pi ON p.id = pi.portfolio_id
		WHERE pi.investment_id = ?
	`

	var id string
	var userID string
	var createdAt string
	var updatedAt string

	err := r.db.QueryRowContext(ctx, query, investmentID.Value).Scan(
		&id,
		&userID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	portfolio := domain.NewPortfolio(domain.NewPortfolioID(id), userID)

	// Load investments
	investments, err := r.loadPortfolioInvestments(ctx, portfolio.ID())
	if err != nil {
		return nil, err
	}

	for _, investment := range investments {
		if err := portfolio.AddInvestment(investment); err != nil {
			return nil, err
		}
	}

	return portfolio, nil
}
