package app

import (
	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app/domain"
)

type EmployeeRepository interface {
	FindByID(id uuid.UUID) (*domain.Employee, error)
	Insert(employee *domain.Employee) error
}
