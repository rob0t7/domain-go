package domain

import "github.com/google/uuid"

type Employee struct {
	id        uuid.UUID
	firstName string
	lastName  string
	salary    int32
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

func (e *Employee) Salary() int32 {
	return e.salary
}
