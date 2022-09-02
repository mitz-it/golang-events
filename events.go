package events

import (
	"fmt"
	"strings"

	"go.uber.org/dig"
	"golang.org/x/exp/slices"
)

type IDomainEventHandler interface {
	Handle(domainEvent IDomainEvent)
}

type IDomainEvent interface {
}

type IEventHandler interface {
	Handle(IEvent)
}

type IEvent interface {
}

type IEventDispatcher interface {
	AddDomainEvent(IDomainEvent)
	CommitDomainEventsStack()
	DispatchEvent(IEvent)
	AddEvent(IEvent)
	CommitEventsStack()
}

type EventDispatcherParams struct {
	dig.In

	DomainEventHandlers []IDomainEventHandler `group:"DomainEventHandlers"`
	EventHandlers       []IEventHandler       `group:"EventHandlers"`
}

type EventDispatcher struct {
	domainEventHandlers []IDomainEventHandler
	eventHandlers       []IEventHandler

	domainEvents []IDomainEvent
	events       []IEvent
}

func (eventDispatcher *EventDispatcher) AddDomainEvent(event IDomainEvent) {
	eventDispatcher.domainEvents = append(eventDispatcher.domainEvents, event)
}

func (eventDispatcher *EventDispatcher) AddEvent(event IEvent) {
	eventDispatcher.events = append(eventDispatcher.events, event)
}

func (eventDispatcher *EventDispatcher) dispatchDomainEvent(event IDomainEvent) {

	position := slices.IndexFunc(eventDispatcher.domainEventHandlers, func(handler IDomainEventHandler) bool {
		handlerName := fmt.Sprintf("%T", handler)
		eventName := fmt.Sprintf("%T", event)
		return strings.Contains(handlerName, eventName)
	})

	eventDispatcher.domainEventHandlers[position].Handle(event)
}

func (eventDispatcher *EventDispatcher) DispatchEvent(event IEvent) {
	position := slices.IndexFunc(eventDispatcher.eventHandlers, func(handler IEventHandler) bool {
		handlerName := fmt.Sprintf("%T", handler)
		eventName := fmt.Sprintf("%T", event)
		return strings.Contains(handlerName, eventName)
	})

	eventDispatcher.eventHandlers[position].Handle(event)
}

func (eventDispatcher *EventDispatcher) CommitDomainEventsStack() {
	for _, event := range eventDispatcher.domainEvents {
		eventDispatcher.dispatchDomainEvent(event)
	}
}

func (eventDispatcher *EventDispatcher) CommitEventsStack() {
	for _, event := range eventDispatcher.events {
		eventDispatcher.DispatchEvent(event)
	}
}

func NewEventDispatcher(params EventDispatcherParams) IEventDispatcher {
	return &EventDispatcher{domainEventHandlers: params.DomainEventHandlers, eventHandlers: params.EventHandlers}
}
