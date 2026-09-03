package domain

import (
	"fmt"
	"time"

	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

// Book is the aggregate root of the catalog bounded context.
// Outside code mutates catalog state only through Book methods.
type Book struct {
	isbn   ISBN
	title  string
	author string
	copies Copies
	events []shared.DomainEvent
}

func RegisterBook(isbn ISBN, title, author string, copies int, now time.Time) (*Book, error) {
	if isbn.IsZero() {
		return nil, fmt.Errorf("book: isbn is required")
	}
	if n := len([]rune(title)); n < 1 || n > 200 {
		return nil, fmt.Errorf("book: invalid title")
	}
	if n := len([]rune(author)); n < 1 || n > 120 {
		return nil, fmt.Errorf("book: invalid author")
	}
	c, err := NewCopies(copies)
	if err != nil {
		return nil, err
	}
	b := &Book{isbn: isbn, title: title, author: author, copies: c}
	b.raise(BookRegistered{ISBN: isbn, Title: title, OccurredAt: now.UTC()})
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

func (b *Book) DomainEvents() []shared.DomainEvent {
	out := make([]shared.DomainEvent, len(b.events))
	copy(out, b.events)
	return out
}

func (b *Book) PullEvents() []shared.DomainEvent {
	out := b.DomainEvents()
	b.events = nil
	return out
}

func (b *Book) raise(e shared.DomainEvent) {
	b.events = append(b.events, e)
}
