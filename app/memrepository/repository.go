package memrepository

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/domain"
)

type Repository map[uuid.UUID]*domain.Company

func New() Repository {
	return make(map[uuid.UUID]*domain.Company)
}

func (r Repository) FindAll() []*domain.Company {
	var companies []*domain.Company
	for _, c := range r {
		company := *c
		companies = append(companies, &company)
	}
	sort.SliceStable(companies, func(i, j int) bool {
		return companies[i].Name() < companies[j].Name()
	})
	return companies
}

func (r Repository) FindByID(id uuid.UUID) (*domain.Company, error) {
	company, found := r[id]
	if !found {
		return nil, app.ErrNotFound
	}
	out := *company
	return &out, nil
}

func (r Repository) Insert(company *domain.Company) error {
	if _, found := r[company.ID()]; found {
		return fmt.Errorf("%w: Company with ID = \"%s\" already exists", app.ErrConflict, company.ID())
	}
	value := *company
	r[company.ID()] = &value
	return nil
}

func (r Repository) Update(company *domain.Company) error {
	if _, found := r[company.ID()]; !found {
		return app.ErrNotFound
	}
	newCopy := *company
	r[company.ID()] = &newCopy
	return nil
}

func (r Repository) Delete(company *domain.Company) error {
	if _, found := r[company.ID()]; !found {
		return app.ErrNotFound
	}
	delete(r, company.ID())
	return nil
}

func (r Repository) Reset() {
	for k := range r {
		delete(r, k)
	}
}
