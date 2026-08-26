package main

import (
	"fmt"
	"time"

	catcmd "github.com/cagatayozata/ddd-library-lending/internal/catalog/application/command"
	"github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
	catmem "github.com/cagatayozata/ddd-library-lending/internal/catalog/infrastructure/memory"
	lendcmd "github.com/cagatayozata/ddd-library-lending/internal/lending/application/command"
	lending "github.com/cagatayozata/ddd-library-lending/internal/lending/domain"
	lendmem "github.com/cagatayozata/ddd-library-lending/internal/lending/infrastructure/memory"
	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

func main() {
	now := time.Date(2026, 8, 16, 10, 0, 0, 0, time.UTC)
	books := catmem.NewBookRepository()
	loans := lendmem.NewLoanRepository()

	isbn, err := domain.NewISBN("9780132350884")
	must(err)

	book, registered, err := catcmd.RegisterBookHandler{Books: books}.Handle(catcmd.RegisterBookCommand{
		ISBN: isbn, Title: "Clean Code", Author: "Robert C. Martin", Copies: 3,
	}, now)
	must(err)

	fmt.Println("=== catalog ===")
	fmt.Printf("%q by %s | copies=%d | events=%v\n",
		book.Title(), book.Author(), book.Copies().Available(), names(registered))

	member, _ := lending.NewMemberID("M-42")
	loan, borrowed, err := lendcmd.BorrowBookHandler{Books: books, Loans: loans}.Handle(lendcmd.BorrowBookCommand{
		LoanID: "L-100", MemberID: member, ISBN: isbn,
	}, now)
	must(err)

	fmt.Println("\n=== lending (refs ISBN only) ===")
	fmt.Printf("loan %s | member=%s | isbn=%s | due=%s | status=%s | events=%v\n",
		loan.ID(), loan.MemberID(), loan.ISBN(), loan.DueDate(), loan.Status(), names(borrowed))

	renewed, err := lendcmd.RenewLoanHandler{Loans: loans}.Handle("L-100", now.Add(24*time.Hour))
	must(err)
	fmt.Printf("renewed | due=%s | events=%v\n", loan.DueDate(), names(renewed))

	returned, err := lendcmd.ReturnBookHandler{Books: books, Loans: loans}.Handle("L-100", now.Add(48*time.Hour))
	must(err)
	stored, _ := books.FindByISBN(isbn)
	fmt.Printf("returned | events=%v | catalog copies=%d\n", names(returned), stored.Copies().Available())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func names(evts []shared.DomainEvent) []string {
	out := make([]string, len(evts))
	for i, e := range evts {
		out[i] = e.EventName()
	}
	return out
}
