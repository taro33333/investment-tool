package service

import (
	"encoding/json"
	"moneyget/internal/domain"
	"time"
)

type EventStore struct {
	db EventStoreDB
}

type EventStoreDB interface {
	Store(event domain.DomainEvent) error
}

type StoredEvent struct {
	ID            string          `json:"id"`
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	EventType     string          `json:"event_type"`
	EventData     json.RawMessage `json:"event_data"`
	OccurredAt    time.Time       `json:"occurred_at"`
}

func NewEventStore(db EventStoreDB) *EventStore {
	return &EventStore{
		db: db,
	}
}

func (s *EventStore) SaveEvent(event domain.DomainEvent) error {
	return s.db.Store(event)
}

func GetEventType(event domain.DomainEvent) string {
	switch event.(type) {
	case domain.InvestmentCreatedEvent:
		return "InvestmentCreated"
	case domain.PortfolioUpdatedEvent:
		return "PortfolioUpdated"
	default:
		return "Unknown"
	}
}
