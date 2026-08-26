package domain

type LoanRepository interface {
	Save(loan *Loan) error
	FindByID(id LoanID) (*Loan, error)
}
