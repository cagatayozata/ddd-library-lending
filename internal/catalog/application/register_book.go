package application

import (
	"time"

	"github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
)

// RegisterBook is a thin proxy: build aggregate root → save.
// Controllers can call this, or call the domain + repository directly — DDD does not require CQRS.
type RegisterBook struct {
	Books domain.BookRepository
}

type RegisterBookRequest struct {
	ISBN   domain.ISBN
	Title  string
	Author string
	Copies int
}

func (s RegisterBook) Handle(req RegisterBookRequest, now time.Time) (*domain.Book, error) {
	book, err := domain.RegisterBook(req.ISBN, req.Title, req.Author, req.Copies, now)
	if err != nil {
		return nil, err
	}
	if err := s.Books.Save(book); err != nil {
		return nil, err
	}
	return book, nil
}
