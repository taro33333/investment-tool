package sqlite

import (
	"database/sql"
	"encoding/json"
	"moneyget/internal/domain"
)

type EventStoreDB struct {
	db *sql.DB
}

func NewEventStoreDB(db *sql.DB) *EventStoreDB {
	return &EventStoreDB{db: db}
}

func (e *EventStoreDB) Store(event domain.DomainEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO events (event_type, event_data, occurred_at)
		VALUES (?, ?, ?)
	`

	_, err = e.db.Exec(query, getEventType(event), eventData, event.OccurredAt())
	return err
}

func getEventType(event domain.DomainEvent) string {
	switch event.(type) {
	case domain.InvestmentCreatedEvent:
		return "InvestmentCreated"
	case domain.PortfolioUpdatedEvent:
		return "PortfolioUpdated"
	default:
		return "Unknown"
	}
}
