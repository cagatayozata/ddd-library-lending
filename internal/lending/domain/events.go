package domain

import "time"

type LoanBorrowed struct {
	LoanID     LoanID
	MemberID   MemberID
	ISBN       ISBN
	OccurredAt time.Time
}

func (e LoanBorrowed) EventName() string { return "LoanBorrowed" }

type LoanReturnedEvent struct {
	LoanID     LoanID
	OccurredAt time.Time
}

func (e LoanReturnedEvent) EventName() string { return "LoanReturned" }

type LoanRenewed struct {
	LoanID     LoanID
	DueDate    DueDate
	OccurredAt time.Time
}

func (e LoanRenewed) EventName() string { return "LoanRenewed" }
