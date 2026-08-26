package memory

import (
	"fmt"
	"sync"

	"github.com/cagatayozata/ddd-library-lending/internal/lending/domain"
)

type LoanRepository struct {
	mu   sync.RWMutex
	byID map[domain.LoanID]*domain.Loan
}

func NewLoanRepository() *LoanRepository {
	return &LoanRepository{byID: make(map[domain.LoanID]*domain.Loan)}
}

func (r *LoanRepository) Save(loan *domain.Loan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[loan.ID()] = loan
	return nil
}

func (r *LoanRepository) FindByID(id domain.LoanID) (*domain.Loan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	loan, ok := r.byID[id]
	if !ok {
		return nil, fmt.Errorf("loan not found: %s", id)
	}
	return loan, nil
}
