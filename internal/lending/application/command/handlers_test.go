package command_test

import (
	"testing"
	"time"

	catcmd "github.com/cagatayozata/ddd-library-lending/internal/catalog/application/command"
	catmem "github.com/cagatayozata/ddd-library-lending/internal/catalog/infrastructure/memory"
	"github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
	lendcmd "github.com/cagatayozata/ddd-library-lending/internal/lending/application/command"
	lending "github.com/cagatayozata/ddd-library-lending/internal/lending/domain"
	lendmem "github.com/cagatayozata/ddd-library-lending/internal/lending/infrastructure/memory"
)

func TestBorrowAndReturn_Flow(t *testing.T) {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	books := catmem.NewBookRepository()
	loans := lendmem.NewLoanRepository()

	isbn, err := domain.NewISBN("9780132350884")
	if err != nil {
		t.Fatal(err)
	}

	register := catcmd.RegisterBookHandler{Books: books}
	_, _, err = register.Handle(catcmd.RegisterBookCommand{
		ISBN: isbn, Title: "Clean Code", Author: "Robert C. Martin", Copies: 2,
	}, now)
	if err != nil {
		t.Fatal(err)
	}

	member, _ := lending.NewMemberID("M-1")
	borrow := lendcmd.BorrowBookHandler{Books: books, Loans: loans}
	loan, evts, err := borrow.Handle(lendcmd.BorrowBookCommand{
		LoanID: "L-1", MemberID: member, ISBN: isbn,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if loan.Status() != lending.LoanActive {
		t.Fatalf("status=%s", loan.Status())
	}
	if len(evts) != 1 || evts[0].EventName() != "LoanBorrowed" {
		t.Fatalf("events=%#v", evts)
	}

	book, _ := books.FindByISBN(isbn)
	if book.Copies().Available() != 1 {
		t.Fatalf("available=%d", book.Copies().Available())
	}

	ret := lendcmd.ReturnBookHandler{Books: books, Loans: loans}
	if _, err = ret.Handle("L-1", now.Add(24*time.Hour)); err != nil {
		t.Fatal(err)
	}
	book, _ = books.FindByISBN(isbn)
	if book.Copies().Available() != 2 {
		t.Fatalf("available=%d", book.Copies().Available())
	}
}
