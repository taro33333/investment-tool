package domain

import (
	"time"
)

type DomainEvent interface {
	OccurredAt() time.Time
}

type InvestmentCreatedEvent struct {
	investmentID InvestmentID
	amount       Money
	occurredAt   time.Time
}

func NewInvestmentCreatedEvent(investmentID InvestmentID, amount Money) InvestmentCreatedEvent {
	return InvestmentCreatedEvent{
		investmentID: investmentID,
		amount:       amount,
		occurredAt:   time.Now(),
	}
}

func (e InvestmentCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

type PortfolioUpdatedEvent struct {
	portfolioID PortfolioID
	totalAmount Money
	occurredAt  time.Time
}

func NewPortfolioUpdatedEvent(portfolioID PortfolioID, totalAmount Money) PortfolioUpdatedEvent {
	return PortfolioUpdatedEvent{
		portfolioID: portfolioID,
		totalAmount: totalAmount,
		occurredAt:  time.Now(),
	}
}

func (e PortfolioUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

type DomainEventPublisher interface {
	Publish(event DomainEvent) error
	Subscribe(handler func(DomainEvent)) error
}
