package domain

import "errors"

type Money struct {
	Amount   float64
	Currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, ErrInvalidInvestmentAmount
	}
	if currency == "" {
		return Money{}, errors.New("currency is required")
	}
	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return NewMoney(m.Amount+other.Amount, m.Currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}
	return NewMoney(m.Amount-other.Amount, m.Currency)
}

func (m Money) Multiply(factor float64) (Money, error) {
	return NewMoney(m.Amount*factor, m.Currency)
}

func (m Money) IsZero() bool {
	return m.Amount == 0
}

func (m Money) IsGreaterThan(other Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount > other.Amount
}

func (m Money) IsLessThan(other Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount < other.Amount
}
