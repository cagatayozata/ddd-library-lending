package command

import (
	"time"

	"github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
	lending "github.com/cagatayozata/ddd-library-lending/internal/lending/domain"
	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

// BorrowBookCommand coordinates catalog + lending by ISBN.
type BorrowBookCommand struct {
	LoanID   lending.LoanID
	MemberID lending.MemberID
	ISBN     domain.ISBN
}

type BorrowBookHandler struct {
	Books domain.BookRepository
	Loans lending.LoanRepository
}

func (h BorrowBookHandler) Handle(cmd BorrowBookCommand, now time.Time) (*lending.Loan, []shared.DomainEvent, error) {
	book, err := h.Books.FindByISBN(cmd.ISBN)
	if err != nil {
		return nil, nil, err
	}
	if err := book.ReserveCopy(); err != nil {
		return nil, nil, err
	}

	loanISBN, err := lending.NewISBN(cmd.ISBN.String())
	if err != nil {
		return nil, nil, err
	}
	loan, err := lending.Borrow(cmd.LoanID, cmd.MemberID, loanISBN, now)
	if err != nil {
		return nil, nil, err
	}

	if err := h.Books.Save(book); err != nil {
		return nil, nil, err
	}
	if err := h.Loans.Save(loan); err != nil {
		return nil, nil, err
	}
	return loan, loan.DrainDomainEvents(), nil
}

type ReturnBookHandler struct {
	Books domain.BookRepository
	Loans lending.LoanRepository
}

func (h ReturnBookHandler) Handle(loanID lending.LoanID, now time.Time) ([]shared.DomainEvent, error) {
	loan, err := h.Loans.FindByID(loanID)
	if err != nil {
		return nil, err
	}
	isbn, err := domain.NewISBN(loan.ISBN().String())
	if err != nil {
		return nil, err
	}
	book, err := h.Books.FindByISBN(isbn)
	if err != nil {
		return nil, err
	}

	if err := loan.Return(now); err != nil {
		return nil, err
	}
	book.ReleaseCopy()

	if err := h.Loans.Save(loan); err != nil {
		return nil, err
	}
	if err := h.Books.Save(book); err != nil {
		return nil, err
	}
	return loan.DrainDomainEvents(), nil
}

type RenewLoanHandler struct {
	Loans lending.LoanRepository
}

func (h RenewLoanHandler) Handle(loanID lending.LoanID, now time.Time) ([]shared.DomainEvent, error) {
	loan, err := h.Loans.FindByID(loanID)
	if err != nil {
		return nil, err
	}
	if err := loan.Renew(now); err != nil {
		return nil, err
	}
	if err := h.Loans.Save(loan); err != nil {
		return nil, err
	}
	return loan.DrainDomainEvents(), nil
}
