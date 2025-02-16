package service

import (
	"moneyget/internal/domain"
	"sync"
)

type EventDispatcher struct {
	handlers []func(domain.DomainEvent)
	mu       sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make([]func(domain.DomainEvent), 0),
	}
}

func (d *EventDispatcher) Subscribe(handler func(domain.DomainEvent)) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, handler)
	return nil
}

func (d *EventDispatcher) Publish(event domain.DomainEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, handler := range d.handlers {
		handler(event)
	}
	return nil
}
