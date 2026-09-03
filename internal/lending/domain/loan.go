package domain

import (
	"fmt"
	"time"

	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

type LoanStatus string

const (
	LoanActive   LoanStatus = "active"
	LoanReturned LoanStatus = "returned"
)

const (
	DefaultLoanDays = 14
	RenewExtraDays  = 7
	MaxRenewals     = 2
)

// Loan is the aggregate root of the lending bounded context.
// Rules (overdue, renew limits, return once) live here — not in application handlers.
type Loan struct {
	id         LoanID
	memberID   MemberID
	isbn       ISBN
	dueDate    DueDate
	status     LoanStatus
	renewCount int
	events     []shared.DomainEvent
}

func Borrow(id LoanID, member MemberID, isbn ISBN, at time.Time) (*Loan, error) {
	if id.IsZero() {
		return nil, fmt.Errorf("loan: id is required")
	}
	if member.IsZero() {
		return nil, fmt.Errorf("loan: member is required")
	}
	if isbn.IsZero() {
		return nil, fmt.Errorf("loan: isbn is required")
	}
	loan := &Loan{
		id:       id,
		memberID: member,
		isbn:     isbn,
		dueDate:  DueDateFromDays(at, DefaultLoanDays),
		status:   LoanActive,
	}
	loan.raise(LoanBorrowed{LoanID: id, MemberID: member, ISBN: isbn, OccurredAt: at.UTC()})
	return loan, nil
}

func (l *Loan) ID() LoanID         { return l.id }
func (l *Loan) MemberID() MemberID { return l.memberID }
func (l *Loan) ISBN() ISBN         { return l.isbn }
func (l *Loan) DueDate() DueDate   { return l.dueDate }
func (l *Loan) Status() LoanStatus { return l.status }
func (l *Loan) RenewCount() int    { return l.renewCount }

func (l *Loan) IsOverdue(now time.Time) bool {
	return l.status == LoanActive && l.dueDate.IsBefore(now)
}

func (l *Loan) Return(at time.Time) error {
	if l.status == LoanReturned {
		return fmt.Errorf("loan: already returned")
	}
	l.status = LoanReturned
	l.raise(LoanReturnedEvent{LoanID: l.id, OccurredAt: at.UTC()})
	return nil
}

func (l *Loan) Renew(now time.Time) error {
	if l.status != LoanActive {
		return fmt.Errorf("loan: cannot renew")
	}
	if l.IsOverdue(now) {
		return fmt.Errorf("loan: overdue")
	}
	if l.renewCount >= MaxRenewals {
		return fmt.Errorf("loan: max renewals reached")
	}
	l.dueDate = l.dueDate.AddDays(RenewExtraDays)
	l.renewCount++
	l.raise(LoanRenewed{LoanID: l.id, DueDate: l.dueDate, OccurredAt: now.UTC()})
	return nil
}

func (l *Loan) DomainEvents() []shared.DomainEvent {
	out := make([]shared.DomainEvent, len(l.events))
	copy(out, l.events)
	return out
}

func (l *Loan) PullEvents() []shared.DomainEvent {
	out := l.DomainEvents()
	l.events = nil
	return out
}

func (l *Loan) raise(e shared.DomainEvent) {
	l.events = append(l.events, e)
}
