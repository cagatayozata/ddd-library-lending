package memory

import (
	"fmt"
	"sync"

	"github.com/cagatayozata/ddd-library-lending/internal/catalog/domain"
)

type BookRepository struct {
	mu   sync.RWMutex
	byID map[string]*domain.Book
}

func NewBookRepository() *BookRepository {
	return &BookRepository{byID: make(map[string]*domain.Book)}
}

func (r *BookRepository) Save(book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[book.ISBN().String()] = book
	return nil
}

func (r *BookRepository) FindByISBN(isbn domain.ISBN) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.byID[isbn.String()]
	if !ok {
		return nil, fmt.Errorf("book not found: %s", isbn)
	}
	return b, nil
}
