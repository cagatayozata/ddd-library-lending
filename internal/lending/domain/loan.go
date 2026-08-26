package domain

import (
	"fmt"
	"strings"
	"time"

	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

// ISBN mirrors catalog identity — lending never imports Book, only this value.
type ISBN struct {
	value string
}

func NewISBN(raw string) (ISBN, error) {
	v := strings.ReplaceAll(strings.TrimSpace(raw), "-", "")
	if len(v) != 10 && len(v) != 13 {
		return ISBN{}, fmt.Errorf("isbn: must be 10 or 13 digits")
	}
	return ISBN{value: v}, nil
}

func (i ISBN) String() string { return i.value }

type MemberID string

func NewMemberID(raw string) (MemberID, error) {
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("member id: empty")
	}
	return MemberID(raw), nil
}

func (m MemberID) String() string { return string(m) }

type LoanID string

func NewLoanID(raw string) (LoanID, error) {
	if raw == "" {
		return "", fmt.Errorf("loan id: empty")
	}
	return LoanID(raw), nil
}

func (id LoanID) String() string { return string(id) }

// DueDate is a Value Object (date-only semantics).
type DueDate struct {
	day time.Time
}

func NewDueDate(t time.Time) DueDate {
	utc := t.UTC()
	return DueDate{day: time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)}
}

func DueDateFromDays(from time.Time, days int) DueDate {
	return NewDueDate(from.AddDate(0, 0, days))
}

func (d DueDate) AddDays(days int) DueDate { return NewDueDate(d.day.AddDate(0, 0, days)) }
func (d DueDate) IsBefore(t time.Time) bool {
	return d.day.Before(NewDueDate(t).day)
}
func (d DueDate) Equals(other DueDate) bool { return d.day.Equal(other.day) }
func (d DueDate) String() string            { return d.day.Format("2006-01-02") }

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

// Loan is the Aggregate Root of the lending bounded context.
type Loan struct {
	shared.AggregateRoot

	id         LoanID
	memberID   MemberID
	isbn       ISBN
	dueDate    DueDate
	status     LoanStatus
	renewCount int
}

func Borrow(id LoanID, member MemberID, isbn ISBN, at time.Time) (*Loan, error) {
	if err := shared.RequirePresent(id, "loan.id"); err != nil {
		return nil, err
	}
	if err := shared.RequirePresent(member, "member.id"); err != nil {
		return nil, err
	}
	if isbn.value == "" {
		return nil, fmt.Errorf("invariant: isbn is required")
	}
	loan := &Loan{
		id:       id,
		memberID: member,
		isbn:     isbn,
		dueDate:  DueDateFromDays(at, DefaultLoanDays),
		status:   LoanActive,
	}
	loan.Raise(LoanBorrowed{LoanID: id, MemberID: member, ISBN: isbn, OccurredAt: at.UTC()})
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
	l.Raise(LoanReturnedEvent{LoanID: l.id, OccurredAt: at.UTC()})
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
	l.Raise(LoanRenewed{LoanID: l.id, DueDate: l.dueDate, OccurredAt: now.UTC()})
	return nil
}
