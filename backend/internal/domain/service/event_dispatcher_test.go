package service

import (
	"moneyget/internal/domain"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	data       string
	occurredAt time.Time
}

func (e testEvent) Name() string {
	return "TestEvent"
}

func (e testEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func TestEventDispatcher(t *testing.T) {
	dispatcher := NewEventDispatcher()

	t.Run("Subscribe and Publish", func(t *testing.T) {
		receivedEvents := make([]domain.DomainEvent, 0)
		var mu sync.Mutex

		err := dispatcher.Subscribe(func(event domain.DomainEvent) {
			mu.Lock()
			receivedEvents = append(receivedEvents, event)
			mu.Unlock()
		})
		assert.NoError(t, err)

		testEvent1 := testEvent{
			data:       "test1",
			occurredAt: time.Now(),
		}
		testEvent2 := testEvent{
			data:       "test2",
			occurredAt: time.Now(),
		}

		err = dispatcher.Publish(testEvent1)
		assert.NoError(t, err)
		err = dispatcher.Publish(testEvent2)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(receivedEvents))
		assert.Equal(t, testEvent1, receivedEvents[0])
		assert.Equal(t, testEvent2, receivedEvents[1])
	})

	t.Run("Multiple Subscribers", func(t *testing.T) {
		counter := 0
		var mu sync.Mutex

		for i := 0; i < 3; i++ {
			err := dispatcher.Subscribe(func(event domain.DomainEvent) {
				mu.Lock()
				counter++
				mu.Unlock()
			})
			assert.NoError(t, err)
		}

		event := testEvent{
			data:       "test",
			occurredAt: time.Now(),
		}
		err := dispatcher.Publish(event)
		assert.NoError(t, err)

		assert.Equal(t, 3, counter)
	})
}
