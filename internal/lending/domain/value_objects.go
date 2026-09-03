package domain

import (
	"fmt"
	"strings"
	"time"
)

// LoanID is a Value Object.
type LoanID struct {
	value string
}

func NewLoanID(raw string) (LoanID, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return LoanID{}, fmt.Errorf("loan id: empty")
	}
	return LoanID{value: v}, nil
}

func (id LoanID) String() string           { return id.value }
func (id LoanID) Equals(other LoanID) bool { return id.value == other.value }
func (id LoanID) IsZero() bool             { return id.value == "" }

// MemberID is a Value Object.
type MemberID struct {
	value string
}

func NewMemberID(raw string) (MemberID, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return MemberID{}, fmt.Errorf("member id: empty")
	}
	return MemberID{value: v}, nil
}

func (m MemberID) String() string              { return m.value }
func (m MemberID) Equals(other MemberID) bool  { return m.value == other.value }
func (m MemberID) IsZero() bool                { return m.value == "" }

// ISBN is a Value Object (lending-side identity for a book — not a Book entity).
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

func (i ISBN) String() string         { return i.value }
func (i ISBN) Equals(other ISBN) bool { return i.value == other.value }
func (i ISBN) IsZero() bool           { return i.value == "" }

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
