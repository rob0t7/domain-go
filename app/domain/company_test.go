package domain_test

import (
	"testing"

	"github.com/rob0t7/domain-go/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCompanySuccessfully(t *testing.T) {
	companyName := "ACME INC"
	company, event := domain.RegisterCompany(companyName)
	assert.Equal(t, companyName, company.Name())
	assert.Equal(t, companyName, event.CompanyName)
	assert.Equal(t, company.ID(), event.CompanyID)
}

func TestFailedCompanyRegistration(t *testing.T) {
	companyName := ""
	assert.Panics(t, func() {
		domain.RegisterCompany(companyName)
	})
}

func TestChangeCompanyName(t *testing.T) {
	company, _ := domain.RegisterCompany("ACME INC")
	newName := "ACME CANADA INC"
	event := company.Rename(newName)
	assert.Equal(t, newName, company.Name())
	assert.Equal(
		t,
		domain.CompanyNameChangedEvent{
			CompanyID:   company.ID(),
			CompanyName: newName,
		},
		event,
	)
}

func TestHireEmployee(t *testing.T) {
	company, _ := domain.RegisterCompany("ACME INC")
	employeeFirstName := "John"
	employeeLastName := "Smith"

	employee, event := company.HireEmployee(employeeFirstName, employeeLastName)
	assert.Len(t, company.Employees(), 1)
	assert.Equal(t, employeeFirstName, employee.FirstName())
	assert.Equal(t, employeeLastName, employee.LastName())
	assert.Equal(
		t,
		domain.EmployeeHiredEvent{
			CompanyID:         company.ID(),
			CompanyName:       company.Name(),
			EmployeeID:        employee.ID(),
			EmployeeFirstName: employeeFirstName,
			EmployeeLastName:  employeeLastName,
		},
		event,
	)
}

func TestFireEmployee(t *testing.T) {
	company, _ := domain.RegisterCompany("ACME INC")
	employee, _ := company.HireEmployee("John", "Smith")

	event := company.FireEmployee(employee.ID())

	assert.Len(t, company.Employees(), 0)
	assert.Equal(
		t,
		domain.EmployeeFiredEvent{
			CompanyID:         company.ID(),
			CompanyName:       company.Name(),
			EmployeeID:        employee.ID(),
			EmployeeFirstName: employee.FirstName(),
			EmployeeLastName:  employee.LastName(),
		},
		event,
	)
}
