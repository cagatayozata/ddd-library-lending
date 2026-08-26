package domain

// DomainEvent is something meaningful that happened in the domain.
// Aggregates raise events; application/infrastructure drain and react (audit, etc.).
type DomainEvent interface {
	EventName() string
}
