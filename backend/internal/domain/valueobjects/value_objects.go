package valueobjects

import (
	"errors"
	"regexp"
)

type Currency struct {
	code string
}

var supportedCurrencies = map[string]bool{
	"JPY": true,
	"USD": true,
	"EUR": true,
}

func NewCurrency(code string) (Currency, error) {
	if !supportedCurrencies[code] {
		return Currency{}, errors.New("unsupported currency code")
	}
	return Currency{code: code}, nil
}

func (c Currency) String() string {
	return c.code
}

type Money struct {
	amount   float64
	currency Currency
}

func NewMoney(amount float64, currencyCode string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}

	currency, err := NewCurrency(currencyCode)
	if err != nil {
		return Money{}, err
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency.code != other.currency.code {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return NewMoney(m.amount+other.amount, m.currency.code)
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency.code != other.currency.code {
		return Money{}, errors.New("cannot subtract money with different currencies")
	}
	return NewMoney(m.amount-other.amount, m.currency.code)
}

func (m Money) Multiply(factor float64) (Money, error) {
	return NewMoney(m.amount*factor, m.currency.code)
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}

type Email struct {
	address string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(address string) (Email, error) {
	if !emailRegex.MatchString(address) {
		return Email{}, errors.New("invalid email format")
	}
	return Email{address: address}, nil
}

func (e Email) String() string {
	return e.address
}

type Password struct {
	hash string
}

func NewPassword(hash string) Password {
	return Password{hash: hash}
}

func (p Password) String() string {
	return p.hash
}
