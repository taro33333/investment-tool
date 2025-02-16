package domain

import (
	"errors"
	"time"
)

type InvestmentID struct {
	Value string // エクスポート
}

func NewInvestmentID(id string) InvestmentID {
	return InvestmentID{Value: id}
}

type InvestmentType string

const (
	Stock      InvestmentType = "STOCK"
	Bond       InvestmentType = "BOND"
	RealEstate InvestmentType = "REAL_ESTATE"
)

type InvestmentStrategy string

const (
	Conservative InvestmentStrategy = "CONSERVATIVE"
	Moderate     InvestmentStrategy = "MODERATE"
	Aggressive   InvestmentStrategy = "AGGRESSIVE"
)

type Investment struct {
	id        InvestmentID
	amount    Money
	typeVal   InvestmentType
	strategy  InvestmentStrategy
	CreatedAt time.Time // エクスポート
	UpdatedAt time.Time // エクスポート
}

func NewInvestment(id InvestmentID, amount Money, typeVal InvestmentType, strategy InvestmentStrategy) (*Investment, error) {
	if !isValidInvestmentType(typeVal) {
		return nil, errors.New("invalid investment type")
	}
	if !isValidInvestmentStrategy(strategy) {
		return nil, errors.New("invalid investment strategy")
	}
	now := time.Now()
	return &Investment{
		id:        id,
		amount:    amount,
		typeVal:   typeVal,
		strategy:  strategy,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (i *Investment) ID() InvestmentID {
	return i.id
}

func (i *Investment) Amount() Money {
	return i.amount
}

func (i *Investment) Type() InvestmentType {
	return i.typeVal
}

func (i *Investment) Strategy() InvestmentStrategy {
	return i.strategy
}

func (i *Investment) UpdateAmount(newAmount Money) error {
	i.amount = newAmount
	i.UpdatedAt = time.Now()
	return nil
}

func isValidInvestmentType(t InvestmentType) bool {
	switch t {
	case Stock, Bond, RealEstate:
		return true
	default:
		return false
	}
}

func isValidInvestmentStrategy(s InvestmentStrategy) bool {
	switch s {
	case Conservative, Moderate, Aggressive:
		return true
	default:
		return false
	}
}
