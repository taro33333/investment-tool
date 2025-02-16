package service

import (
	"moneyget/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// モックEventStoreDB
type mockEventStoreDB struct {
	storedEvents []domain.DomainEvent
}

func (m *mockEventStoreDB) Store(event domain.DomainEvent) error {
	m.storedEvents = append(m.storedEvents, event)
	return nil
}

type testInvestmentCreatedEvent struct {
	investmentID domain.InvestmentID
	createdAt    time.Time
}

func (e testInvestmentCreatedEvent) Name() string {
	return "InvestmentCreated"
}

func (e testInvestmentCreatedEvent) OccurredAt() time.Time {
	return e.createdAt
}

func TestEventStore(t *testing.T) {
	t.Run("SaveEvent success", func(t *testing.T) {
		mockDB := &mockEventStoreDB{storedEvents: make([]domain.DomainEvent, 0)}
		eventStore := NewEventStore(mockDB)

		event := testInvestmentCreatedEvent{
			investmentID: domain.NewInvestmentID("test-id"),
			createdAt:    time.Now(),
		}

		err := eventStore.SaveEvent(event)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(mockDB.storedEvents))
		assert.Equal(t, event, mockDB.storedEvents[0])
	})

	t.Run("GetEventType", func(t *testing.T) {
		money, err := domain.NewMoney(1000, "JPY")
		assert.NoError(t, err)

		// 実際のInvestmentCreatedEventを使用してテスト
		event := domain.NewInvestmentCreatedEvent(
			domain.NewInvestmentID("test-investment-id"),
			money,
		)

		eventType := GetEventType(event)
		assert.Equal(t, "InvestmentCreated", eventType)
	})
}
