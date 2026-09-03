package domain

// DomainEvent is raised by aggregate roots. There is no shared AggregateRoot base type —
// Book and Loan are the roots; they collect their own events.
type DomainEvent interface {
	EventName() string
}
