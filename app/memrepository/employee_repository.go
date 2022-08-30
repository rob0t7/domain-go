package memrepository

import (
	"sort"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/domain"
)

type EmployeeRepository struct {
	employees map[uuid.UUID]*domain.Employee
}

func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{
		employees: make(map[uuid.UUID]*domain.Employee),
	}
}

func (repo *EmployeeRepository) FindAll(companyID uuid.UUID) []*domain.Employee {
	var employees []*domain.Employee

	for _, employee := range repo.employees {
		if employee.CompanyID() == companyID {
			employees = append(employees, domain.NewEmployee(employee.ID(), employee.CompanyID(), employee.FirstName(), employee.LastName(), employee.Salary()))
		}
	}
	sort.SliceStable(employees, func(i, j int) bool {
		if employees[i].LastName() < employees[j].LastName() {
			return true
		}
		if employees[i].LastName() == employees[j].LastName() {
			return employees[i].FirstName() < employees[j].FirstName()
		}
		return false
	})
	return employees
}

func (repo *EmployeeRepository) FindByID(id uuid.UUID) (*domain.Employee, error) {
	employee, found := repo.employees[id]
	if !found {
		return nil, app.ErrNotFound
	}
	return employee, nil
}

func (repo *EmployeeRepository) Reset() {
	repo.employees = make(map[uuid.UUID]*domain.Employee)
}

func (repo *EmployeeRepository) Insert(employee *domain.Employee) error {
	var newEmployee domain.Employee = *employee
	repo.employees[employee.ID()] = &newEmployee
	return nil
}

func (repo *EmployeeRepository) Delete(employee *domain.Employee) error {
	_, found := repo.employees[employee.ID()]
	if !found {
		return app.ErrNotFound
	}
	delete(repo.employees, employee.ID())
	return nil
}
