package domain

import (
	"fmt"
	"strings"
	"time"

	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

// ISBN is a Value Object. Lending references books by ISBN only.
type ISBN struct {
	value string
}

func NewISBN(raw string) (ISBN, error) {
	v := strings.ReplaceAll(strings.TrimSpace(raw), "-", "")
	if len(v) != 10 && len(v) != 13 {
		return ISBN{}, fmt.Errorf("isbn: must be 10 or 13 digits")
	}
	for _, r := range v {
		if r < '0' || r > '9' {
			return ISBN{}, fmt.Errorf("isbn: invalid character")
		}
	}
	return ISBN{value: v}, nil
}

func (i ISBN) String() string         { return i.value }
func (i ISBN) Equals(other ISBN) bool { return i.value == other.value }

// Copies is a Value Object for shelf availability.
type Copies struct {
	available int
}

func NewCopies(n int) (Copies, error) {
	if err := shared.Require(n >= 0, "copies must not be negative"); err != nil {
		return Copies{}, err
	}
	return Copies{available: n}, nil
}

func (c Copies) Available() int { return c.available }

func (c Copies) ReserveOne() (Copies, error) {
	if c.available < 1 {
		return Copies{}, fmt.Errorf("copies: none available")
	}
	return Copies{available: c.available - 1}, nil
}

func (c Copies) ReleaseOne() Copies {
	return Copies{available: c.available + 1}
}

// Book is the Aggregate Root of the catalog bounded context.
type Book struct {
	shared.AggregateRoot

	isbn   ISBN
	title  string
	author string
	copies Copies
}

func RegisterBook(isbn ISBN, title, author string, copies int, now time.Time) (*Book, error) {
	if isbn.value == "" {
		return nil, fmt.Errorf("invariant: isbn is required")
	}
	if err := shared.RequireText(title, "title", 1, 200); err != nil {
		return nil, err
	}
	if err := shared.RequireText(author, "author", 1, 120); err != nil {
		return nil, err
	}
	c, err := NewCopies(copies)
	if err != nil {
		return nil, err
	}
	b := &Book{isbn: isbn, title: title, author: author, copies: c}
	b.Raise(BookRegistered{ISBN: isbn, Title: title, OccurredAt: now.UTC()})
	return b, nil
}

func (b *Book) ISBN() ISBN     { return b.isbn }
func (b *Book) Title() string  { return b.title }
func (b *Book) Author() string { return b.author }
func (b *Book) Copies() Copies { return b.copies }

func (b *Book) ReserveCopy() error {
	next, err := b.copies.ReserveOne()
	if err != nil {
		return err
	}
	b.copies = next
	return nil
}

func (b *Book) ReleaseCopy() {
	b.copies = b.copies.ReleaseOne()
}
