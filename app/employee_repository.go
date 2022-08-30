package app

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EmployeeRepository interface {
	FindAll(companyID uuid.UUID) []*domain.Employee
	FindByID(id uuid.UUID) (*domain.Employee, error)
	Insert(employee *domain.Employee) error
	Delete(employee *domain.Employee) error
}

type TestableEmployeeRepository interface {
	EmployeeRepository
	Reset()
}

func RunEmployeeRepositoryTestSuite(t *testing.T, repository TestableEmployeeRepository, companyRepository TestableCompanyRepository) {
	t.Helper()

	t.Run("searching for a non-existant employee returns ErrNotFound", func(t *testing.T) {
		defer repository.Reset()
		// defer companyRepository.Reset()
		employee, err := repository.FindByID(uuid.New())
		assert.Nil(t, employee)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("persist new employees successfully", func(t *testing.T) {
		defer companyRepository.Reset()
		defer repository.Reset()
		company := domain.New("ACME INC")
		err := companyRepository.Insert(company)
		require.NoError(t, err)
		employee, _ := company.HireEmployee("Jane", "Doe")
		err = repository.Insert(&employee)
		require.NoError(t, err)
		actualEmployee, err := repository.FindByID(employee.ID())
		require.NoError(t, err)
		require.Equal(t, employee, *actualEmployee)
	})

	t.Run("can delete an employee", func(t *testing.T) {
		defer companyRepository.Reset()
		defer repository.Reset()
		company := domain.New("ACME INC")
		err := companyRepository.Insert(company)
		require.NoError(t, err)
		employee, _ := company.HireEmployee("Jane", "Doe")
		err = repository.Insert(&employee)
		require.NoError(t, err)
		err = repository.Delete(&employee)
		require.NoError(t, err)
		err = repository.Delete(&employee)
		require.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("returns a list of employees for a company order by lastname,firstname", func(t *testing.T) {
		defer companyRepository.Reset()
		defer repository.Reset()
		company := domain.New("ACME")
		err := companyRepository.Insert(company)
		require.NoError(t, err)

		employee1, _ := company.HireEmployee("John", "Smith")
		err = repository.Insert(&employee1)
		require.NoError(t, err)
		employee2, _ := company.HireEmployee("Jane", "Doe")
		err = repository.Insert(&employee2)
		require.NoError(t, err)

		actualEmployees := repository.FindAll(company.ID())
		require.Len(t, actualEmployees, 2)
		assert.Equal(t, employee1, *actualEmployees[1])
		assert.Equal(t, employee2, *actualEmployees[0])
	})
}
