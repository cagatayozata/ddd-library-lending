package domain

import (
	"testing"
	"time"
)

func TestBorrow_CreatesActiveLoan(t *testing.T) {
	id, _ := NewLoanID("L-1")
	member, _ := NewMemberID("M-1")
	isbn, _ := NewISBN("9780132350884")
	at := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)

	loan, err := Borrow(id, member, isbn, at)
	if err != nil {
		t.Fatal(err)
	}
	if loan.Status() != LoanActive {
		t.Fatalf("status=%s", loan.Status())
	}
	want := DueDateFromDays(at, DefaultLoanDays)
	if !loan.DueDate().Equals(want) {
		t.Fatalf("due=%s want=%s", loan.DueDate(), want)
	}
	if loan.DomainEvents()[0].EventName() != "LoanBorrowed" {
		t.Fatal("expected LoanBorrowed")
	}
}

func TestReturn_CannotReturnTwice(t *testing.T) {
	loan := mustLoan(t)
	at := time.Date(2026, 8, 20, 12, 0, 0, 0, time.UTC)
	if err := loan.Return(at); err != nil {
		t.Fatal(err)
	}
	if err := loan.Return(at); err == nil {
		t.Fatal("expected already returned")
	}
}

func TestRenew_FailsWhenOverdue(t *testing.T) {
	loan := mustLoan(t)
	now := time.Date(2026, 9, 5, 9, 0, 0, 0, time.UTC)
	if err := loan.Renew(now); err == nil {
		t.Fatal("expected overdue")
	}
}

func TestRenew_ExtendsDueDate(t *testing.T) {
	loan := mustLoan(t)
	now := time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC)
	before := loan.DueDate()
	if err := loan.Renew(now); err != nil {
		t.Fatal(err)
	}
	if !loan.DueDate().Equals(before.AddDays(RenewExtraDays)) {
		t.Fatalf("due=%s", loan.DueDate())
	}
}

func TestLoanID_MemberID_AreValueObjects(t *testing.T) {
	a, _ := NewLoanID("L-1")
	b, _ := NewLoanID("L-1")
	if !a.Equals(b) {
		t.Fatal("loan id equality")
	}
	m1, _ := NewMemberID("M-1")
	m2, _ := NewMemberID("M-1")
	if !m1.Equals(m2) {
		t.Fatal("member id equality")
	}
}

func mustLoan(t *testing.T) *Loan {
	t.Helper()
	id, _ := NewLoanID("L-1")
	member, _ := NewMemberID("M-1")
	isbn, _ := NewISBN("9780132350884")
	at := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	loan, err := Borrow(id, member, isbn, at)
	if err != nil {
		t.Fatal(err)
	}
	return loan
}
