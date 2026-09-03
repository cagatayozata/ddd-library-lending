package domain

import (
	"fmt"
	"strings"
)

// ISBN is a Value Object. Lending references books by ISBN only — never by a Book entity.
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
func (i ISBN) IsZero() bool           { return i.value == "" }

// Copies is a Value Object for shelf availability.
type Copies struct {
	available int
}

func NewCopies(n int) (Copies, error) {
	if n < 0 {
		return Copies{}, fmt.Errorf("copies: must not be negative")
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
