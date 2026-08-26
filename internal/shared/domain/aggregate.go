package domain

// AggregateRoot is the only entry point for changing an aggregate's state.
// Holds identity-related behaviour helpers and collected domain events.
type AggregateRoot struct {
	events []DomainEvent
}

func (a *AggregateRoot) Raise(e DomainEvent) {
	a.events = append(a.events, e)
}

func (a *AggregateRoot) DomainEvents() []DomainEvent {
	out := make([]DomainEvent, len(a.events))
	copy(out, a.events)
	return out
}

// DrainDomainEvents returns and clears pending events (call after successful save).
func (a *AggregateRoot) DrainDomainEvents() []DomainEvent {
	out := a.DomainEvents()
	a.events = nil
	return out
}
