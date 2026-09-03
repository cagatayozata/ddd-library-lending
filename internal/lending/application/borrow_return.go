package application

import (
	"time"

	catalog "github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
	"github.com/cagatayozata/ddd-library-lending/internal/lending/domain"
)

// BorrowBook proxies two aggregate roots. Business rules stay on Book and Loan — not here.
type BorrowBook struct {
	Books catalog.BookRepository
	Loans domain.LoanRepository
}

type BorrowBookRequest struct {
	LoanID   domain.LoanID
	MemberID domain.MemberID
	ISBN     catalog.ISBN
}

func (s BorrowBook) Handle(req BorrowBookRequest, now time.Time) (*domain.Loan, error) {
	book, err := s.Books.FindByISBN(req.ISBN)
	if err != nil {
		return nil, err
	}
	if err := book.ReserveCopy(); err != nil {
		return nil, err
	}
	if err := s.Books.Save(book); err != nil {
		return nil, err
	}

	isbn, err := domain.NewISBN(req.ISBN.String())
	if err != nil {
		return nil, err
	}
	loan, err := domain.Borrow(req.LoanID, req.MemberID, isbn, now)
	if err != nil {
		return nil, err
	}
	if err := s.Loans.Save(loan); err != nil {
		return nil, err
	}
	return loan, nil
}

// RenewLoan: load Loan root → Renew → save.
type RenewLoan struct {
	Loans domain.LoanRepository
}

func (s RenewLoan) Handle(id domain.LoanID, now time.Time) (*domain.Loan, error) {
	loan, err := s.Loans.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := loan.Renew(now); err != nil {
		return nil, err
	}
	if err := s.Loans.Save(loan); err != nil {
		return nil, err
	}
	return loan, nil
}

// ReturnBook: load roots → Return / ReleaseCopy → save.
type ReturnBook struct {
	Books catalog.BookRepository
	Loans domain.LoanRepository
}

func (s ReturnBook) Handle(id domain.LoanID, now time.Time) (*domain.Loan, error) {
	loan, err := s.Loans.FindByID(id)
	if err != nil {
		return nil, err
	}
	isbn, err := catalog.NewISBN(loan.ISBN().String())
	if err != nil {
		return nil, err
	}
	book, err := s.Books.FindByISBN(isbn)
	if err != nil {
		return nil, err
	}

	if err := loan.Return(now); err != nil {
		return nil, err
	}
	book.ReleaseCopy()

	if err := s.Loans.Save(loan); err != nil {
		return nil, err
	}
	if err := s.Books.Save(book); err != nil {
		return nil, err
	}
	return loan, nil
}
