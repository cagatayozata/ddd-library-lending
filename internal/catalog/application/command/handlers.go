package command

import (
	"time"

	"github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
	shared "github.com/cagatayozata/ddd-library-lending/internal/shared/domain"
)

type RegisterBookCommand struct {
	ISBN   domain.ISBN
	Title  string
	Author string
	Copies int
}

type RegisterBookHandler struct {
	Books domain.BookRepository
}

func (h RegisterBookHandler) Handle(cmd RegisterBookCommand, now time.Time) (*domain.Book, []shared.DomainEvent, error) {
	book, err := domain.RegisterBook(cmd.ISBN, cmd.Title, cmd.Author, cmd.Copies, now)
	if err != nil {
		return nil, nil, err
	}
	if err := h.Books.Save(book); err != nil {
		return nil, nil, err
	}
	return book, book.DrainDomainEvents(), nil
}
