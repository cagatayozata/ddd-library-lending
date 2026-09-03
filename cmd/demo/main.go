package main

import (
	"fmt"
	"time"

	catapp "github.com/cagatayozata/ddd-library-lending/internal/catalog/application"
	catalog "github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
	catmem "github.com/cagatayozata/ddd-library-lending/internal/catalog/infrastructure/memory"
	lendapp "github.com/cagatayozata/ddd-library-lending/internal/lending/application"
	lending "github.com/cagatayozata/ddd-library-lending/internal/lending/domain"
	lendmem "github.com/cagatayozata/ddd-library-lending/internal/lending/infrastructure/memory"
	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

func main() {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	books := catmem.NewBookRepository()
	loans := lendmem.NewLoanRepository()

	isbn, err := catalog.NewISBN("9780132350884")
	must(err)

	register := catapp.RegisterBook{Books: books}
	book, err := register.Handle(catapp.RegisterBookRequest{
		ISBN: isbn, Title: "Clean Code", Author: "Robert C. Martin", Copies: 3,
	}, now)
	must(err)

	fmt.Println("=== catalog ===")
	fmt.Printf("%q by %s | copies=%d | events=%v\n",
		book.Title(), book.Author(), book.Copies().Available(), eventNames(book.PullEvents()))

	loanID, _ := lending.NewLoanID("L-100")
	member, _ := lending.NewMemberID("M-42")
	borrow := lendapp.BorrowBook{Books: books, Loans: loans}
	loan, err := borrow.Handle(lendapp.BorrowBookRequest{
		LoanID: loanID, MemberID: member, ISBN: isbn,
	}, now)
	must(err)

	fmt.Println("\n=== lending (via Loan aggregate root) ===")
	fmt.Printf("loan %s | member=%s | isbn=%s | due=%s | status=%s | events=%v\n",
		loan.ID(), loan.MemberID(), loan.ISBN(), loan.DueDate(), loan.Status(), eventNames(loan.PullEvents()))

	renew := lendapp.RenewLoan{Loans: loans}
	loan, err = renew.Handle(loanID, now.Add(24*time.Hour))
	must(err)
	fmt.Printf("renewed | due=%s | events=%v\n", loan.DueDate(), eventNames(loan.PullEvents()))

	ret := lendapp.ReturnBook{Books: books, Loans: loans}
	loan, err = ret.Handle(loanID, now.Add(48*time.Hour))
	must(err)
	stored, _ := books.FindByISBN(isbn)
	fmt.Printf("returned | events=%v | catalog copies=%d\n", eventNames(loan.PullEvents()), stored.Copies().Available())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func eventNames(evts []shared.DomainEvent) []string {
	out := make([]string, len(evts))
	for i, e := range evts {
		out[i] = e.EventName()
	}
	return out
}
