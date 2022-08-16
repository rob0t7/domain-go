package app_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/domain"
)

type MemoryCompanyRepository struct {
	companies map[uuid.UUID]domain.Company
}

func NewMemoryCompanyRepository() *MemoryCompanyRepository {
	return &MemoryCompanyRepository{
		companies: make(map[uuid.UUID]domain.Company),
	}
}

func (r *MemoryCompanyRepository) FindByID(id uuid.UUID) (domain.Company, error) {
	company, found := r.companies[id]
	if !found {
		return domain.Company{}, app.ErrNotFound
	}
	return company, nil
}

func (r *MemoryCompanyRepository) Save(company *domain.Company) error {
	r.companies[company.ID()] = *company
	return nil
}

func TestMemoryCompanyRepository(t *testing.T) {
	repository := NewMemoryCompanyRepository()
	app.CompanyRepositoryTestSuite(t, repository)
}
