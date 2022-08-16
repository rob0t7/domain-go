package app

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app/domain"
	"github.com/stretchr/testify/require"
)

type CompanyRepository interface {
	// FindAll() []domain.Company
	FindByID(id uuid.UUID) (domain.Company, error)
	Save(company *domain.Company) error
	// Delete(company *domain.Company) error
}

func CompanyRepositoryTestSuite(t *testing.T, repository CompanyRepository) {
	t.Run("Save()", func(t *testing.T) {
		company, _ := domain.RegisterCompany("ACME INC")
		require.NoError(t, repository.Save(company))
		c, err := repository.FindByID(company.ID())
		require.NoError(t, err)
		require.Equal(t, *company, c)
	})

	t.Run("FindByID(id)", func(t *testing.T) {
		t.Run("returns ErrNotFound when company does not exist", func(t *testing.T) {
			_, err := repository.FindByID(uuid.New())
			require.ErrorIs(t, err, ErrNotFound)
		})
	})
}
