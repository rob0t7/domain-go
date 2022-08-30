package domain

import "github.com/google/uuid"

type Employee struct {
	id        uuid.UUID
	companyID uuid.UUID
	firstName string
	lastName  string
	salary    int64
}

func NewEmployee(id uuid.UUID, companyID uuid.UUID, firstName, lastName string, salary int64) *Employee {
	var employee Employee
	if id == uuid.Nil {
		panic("id cannot be nil")
	}
	employee.id = id

	if companyID == uuid.Nil {
		panic("companyID cannot be empty")
	}
	employee.companyID = companyID

	if firstName == "" {
		panic("firstName cannot be blank")
	}
	employee.firstName = firstName

	if lastName == "" {
		panic("lastName cannot be blank")
	}
	employee.lastName = lastName

	if salary < 0 {
		panic("salary must be equal or greater than 0")
	}
	employee.salary = salary

	return &employee
}

func (e *Employee) ID() uuid.UUID {
	return e.id
}

func (e *Employee) FirstName() string {
	return e.firstName
}

func (e *Employee) LastName() string {
	return e.lastName
}

func (e *Employee) Salary() int64 {
	return e.salary
}

func (e *Employee) CompanyID() uuid.UUID {
	return e.companyID
}
