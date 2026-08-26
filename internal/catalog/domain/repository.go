package domain

type BookRepository interface {
	Save(book *Book) error
	FindByISBN(isbn ISBN) (*Book, error)
}
