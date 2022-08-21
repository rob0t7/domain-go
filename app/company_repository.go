package app

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ErrConflict = errors.New("confict")

type CompanyRepository interface {
	// FindAll retrieves a collection of Companies sorted alphabetically by name.
	FindAll() []*domain.Company

	// FindByID retreives the company by id. If the company does not exist return ErrNotFound.
	FindByID(id uuid.UUID) (*domain.Company, error)

	// Insert inserts a new Company record into the repository. If the Company cannot be
	// persisted an error is returned.
	Insert(company *domain.Company) error

	// Update the Company in the repository. If the company does not exist throw an ErrNotFound.
	Update(company *domain.Company) error

	// Delete removes a company from the repository.
	Delete(company *domain.Company) error
}

type TestableCompanyRepository interface {
	CompanyRepository
	// Reset causes the Repository to revert to a clean empty state.
	Reset()
}

func RunCompanyRepositoryTestSuite(t *testing.T, repository TestableCompanyRepository) {
	t.Helper()

	t.Run("FindByID() returns ErrNotFound if Company does not exist", func(t *testing.T) {
		defer repository.Reset()
		_, err := repository.FindByID(uuid.New())
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("Can persiste a new Company into the repository", func(t *testing.T) {
		defer repository.Reset()
		expectedCompany := domain.New("ACME INC")
		err := repository.Insert(expectedCompany)
		require.NoError(t, err)

		actualCompany, err := repository.FindByID(expectedCompany.ID())
		require.NoError(t, err)
		require.Equal(t, expectedCompany, actualCompany)
	})

	t.Run("Inserting the Company twice results in a conflict error", func(t *testing.T) {
		defer repository.Reset()
		company := domain.New("ACME INC")
		err := repository.Insert(company)
		require.NoError(t, err)
		err = repository.Insert(company)
		require.ErrorIs(t, err, ErrConflict)
	})

	t.Run("Update throws a ErrNotFound if the Company does not exist", func(t *testing.T) {
		defer repository.Reset()
		company := domain.New("ACME INC")
		err := repository.Update(company)
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("Updates the Company in the repository successfully", func(t *testing.T) {
		defer repository.Reset()
		originalName := "ACME INC"
		updatedName := "Wiley Coyote Enterprises"
		company := domain.New(originalName)
		err := repository.Insert(company)
		require.NoError(t, err)

		company.Rename(updatedName)
		err = repository.Update(company)
		require.NoError(t, err)

		actualCompany, err := repository.FindByID(company.ID())
		require.NoError(t, err)
		require.NotEqual(t, originalName, actualCompany.Name())
		require.Equal(t, updatedName, actualCompany.Name())
	})

	t.Run("Delete() removes a Company from the repository", func(t *testing.T) {
		defer repository.Reset()
		company := domain.New("ACME INCE")
		err := repository.Insert(company)
		require.NoError(t, err)

		err = repository.Delete(company)
		require.NoError(t, err)

		_, err = repository.FindByID(company.ID())
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("Delete() throws an ErrNotFound if the record does not exist", func(t *testing.T) {
		defer repository.Reset()
		company := domain.New("ACME INC")
		err := repository.Delete(company)
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("FindAll retrieves a list of company order by name", func(t *testing.T) {
		defer repository.Reset()
		company1 := domain.New("ACME INC")
		company2 := domain.New("Wiley Coyote Enterprises")
		err := repository.Insert(company2)
		require.NoError(t, err)
		err = repository.Insert(company1)
		require.NoError(t, err)

		companies := repository.FindAll()
		assert.Len(t, companies, 2)
		assert.Equal(t, []*domain.Company{company1, company2}, companies)
	})
}
