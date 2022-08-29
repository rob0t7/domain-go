package memrepository

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/domain"
)

type CompanyRepository struct {
	companies map[uuid.UUID]*domain.Company
}

func New() *CompanyRepository {
	return &CompanyRepository{
		companies: make(map[uuid.UUID]*domain.Company),
	}
}

func (r *CompanyRepository) FindAll() []*domain.Company {
	var companies []*domain.Company
	for _, c := range r.companies {
		company := *c
		companies = append(companies, &company)
	}
	sort.SliceStable(companies, func(i, j int) bool {
		return companies[i].Name() < companies[j].Name()
	})
	return companies
}

func (r *CompanyRepository) FindByID(id uuid.UUID) (*domain.Company, error) {
	company, found := r.companies[id]
	if !found {
		return nil, app.ErrNotFound
	}
	out := *company
	return &out, nil
}

func (r *CompanyRepository) Insert(company *domain.Company) error {
	if _, found := r.companies[company.ID()]; found {
		return fmt.Errorf("%w: Company with ID = \"%s\" already exists", app.ErrConflict, company.ID())
	}
	value := *company
	r.companies[company.ID()] = &value
	return nil
}

func (r *CompanyRepository) Update(company *domain.Company) error {
	if _, found := r.companies[company.ID()]; !found {
		return app.ErrNotFound
	}
	newCopy := *company
	r.companies[company.ID()] = &newCopy
	return nil
}

func (r *CompanyRepository) Delete(company *domain.Company) error {
	if _, found := r.companies[company.ID()]; !found {
		return app.ErrNotFound
	}
	delete(r.companies, company.ID())
	return nil
}

func (r *CompanyRepository) Reset() {
	r.companies = make(map[uuid.UUID]*domain.Company)
}
