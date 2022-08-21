package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type Company struct {
	id        uuid.UUID
	name      string
	employees map[uuid.UUID]Employee
}

func New(name string) *Company {
	id := uuid.New()
	return NewWithID(id, name)
}

func NewWithID(id uuid.UUID, name string) *Company {
	if id == uuid.Nil {
		panic("id cannot be blank")
	}
	if name == "" {
		panic("name cannot be blank")
	}
	return &Company{
		id:        id,
		name:      name,
		employees: make(map[uuid.UUID]Employee),
	}
}

func (c *Company) ID() uuid.UUID {
	return c.id
}

func (c *Company) Name() string {
	return c.name
}

func (c *Company) String() string {
	return fmt.Sprintf(`Company[ID="%s", Name="%s"]`, c.id, c.name)
}

func (c *Company) Employees() []Employee {
	var employees []Employee
	for _, employee := range c.employees {
		employees = append(employees, employee)
	}
	return employees
}

type CompanyRegisteredEvent struct {
	CompanyID   uuid.UUID
	CompanyName string
}

func RegisterCompany(name string) (*Company, CompanyRegisteredEvent) {
	var company Company
	company.id = uuid.New()
	company.employees = make(map[uuid.UUID]Employee)
	company.Rename(name)

	event := CompanyRegisteredEvent{
		CompanyID:   company.id,
		CompanyName: company.Name(),
	}
	return &company, event
}

type CompanyNameChangedEvent struct {
	CompanyID   uuid.UUID
	CompanyName string
}

func (c *Company) Rename(name string) CompanyNameChangedEvent {
	if name == "" {
		panic("company name cannot be blank")
	}
	c.name = name
	return CompanyNameChangedEvent{
		CompanyID:   c.id,
		CompanyName: c.name,
	}
}

type EmployeeHiredEvent struct {
	CompanyID         uuid.UUID
	CompanyName       string
	EmployeeID        uuid.UUID
	EmployeeFirstName string
	EmployeeLastName  string
}

func (c *Company) HireEmployee(firstName, lastName string) (employee Employee, event EmployeeHiredEvent) {
	if firstName == "" {
		panic("employee firstName cannot be blank")
	}
	employee.firstName = firstName

	if lastName == "" {
		panic("employee lastName cannot be blank")
	}
	employee.lastName = lastName

	employee.id = uuid.New()
	c.employees[employee.id] = employee

	return employee, EmployeeHiredEvent{
		CompanyID:         c.id,
		CompanyName:       c.name,
		EmployeeID:        employee.id,
		EmployeeFirstName: employee.firstName,
		EmployeeLastName:  employee.lastName,
	}

}

type EmployeeFiredEvent struct {
	CompanyID         uuid.UUID
	CompanyName       string
	EmployeeID        uuid.UUID
	EmployeeFirstName string
	EmployeeLastName  string
}

func (c *Company) FireEmployee(employeeID uuid.UUID) (event EmployeeFiredEvent) {
	employee := c.employees[employeeID]
	event.CompanyID = c.id
	event.CompanyName = c.name
	event.EmployeeID = employee.id
	event.EmployeeFirstName = employee.firstName
	event.EmployeeLastName = employee.lastName
	delete(c.employees, employee.id)
	return
}
