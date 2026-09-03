package application_test

import (
	"testing"
	"time"

	catapp "github.com/cagatayozata/ddd-library-lending/internal/catalog/application"
	catalog "github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
	catmem "github.com/cagatayozata/ddd-library-lending/internal/catalog/infrastructure/memory"
	lendapp "github.com/cagatayozata/ddd-library-lending/internal/lending/application"
	lending "github.com/cagatayozata/ddd-library-lending/internal/lending/domain"
	lendmem "github.com/cagatayozata/ddd-library-lending/internal/lending/infrastructure/memory"
)

func TestBorrowAndReturn_Flow(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	books := catmem.NewBookRepository()
	loans := lendmem.NewLoanRepository()

	isbn, err := catalog.NewISBN("9780132350884")
	if err != nil {
		t.Fatal(err)
	}

	register := catapp.RegisterBook{Books: books}
	_, err = register.Handle(catapp.RegisterBookRequest{
		ISBN: isbn, Title: "Clean Code", Author: "Robert C. Martin", Copies: 2,
	}, now)
	if err != nil {
		t.Fatal(err)
	}

	loanID, _ := lending.NewLoanID("L-1")
	member, _ := lending.NewMemberID("M-1")
	borrow := lendapp.BorrowBook{Books: books, Loans: loans}
	loan, err := borrow.Handle(lendapp.BorrowBookRequest{
		LoanID: loanID, MemberID: member, ISBN: isbn,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if loan.Status() != lending.LoanActive {
		t.Fatalf("status=%s", loan.Status())
	}

	book, err := books.FindByISBN(isbn)
	if err != nil {
		t.Fatal(err)
	}
	if book.Copies().Available() != 1 {
		t.Fatalf("available=%d", book.Copies().Available())
	}

	ret := lendapp.ReturnBook{Books: books, Loans: loans}
	if _, err = ret.Handle(loanID, now.Add(24*time.Hour)); err != nil {
		t.Fatal(err)
	}
	book, err = books.FindByISBN(isbn)
	if err != nil {
		t.Fatal(err)
	}
	if book.Copies().Available() != 2 {
		t.Fatalf("available=%d", book.Copies().Available())
	}
}
