package domain

import (
	"testing"
	"time"
)

func TestCopies_ReserveAndRelease(t *testing.T) {
	c, err := NewCopies(2)
	if err != nil {
		t.Fatal(err)
	}
	c, err = c.ReserveOne()
	if err != nil {
		t.Fatal(err)
	}
	if c.Available() != 1 {
		t.Fatalf("available=%d", c.Available())
	}
	c, err = c.ReserveOne()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = c.ReserveOne(); err == nil {
		t.Fatal("expected none available")
	}
	c = c.ReleaseOne()
	if c.Available() != 1 {
		t.Fatalf("available=%d", c.Available())
	}
}

func TestISBN_ValueEquality(t *testing.T) {
	a, _ := NewISBN("978-0-13-235088-4")
	b, _ := NewISBN("9780132350884")
	if !a.Equals(b) {
		t.Fatal("normalized ISBNs should be equal")
	}
}

func TestBook_RegisterRaisesEvent(t *testing.T) {
	isbn, _ := NewISBN("9780132350884")
	book, err := RegisterBook(isbn, "Clean Code", "Robert C. Martin", 3, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	evts := book.DomainEvents()
	if len(evts) != 1 || evts[0].EventName() != "BookRegistered" {
		t.Fatalf("events=%#v", evts)
	}
}
