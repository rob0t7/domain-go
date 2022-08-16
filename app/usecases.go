package app

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app/domain"
)

var ErrNotFound = errors.New("not found")

type CompanyService struct {
	repository CompanyRepository
}

func NewCompanyService(repository CompanyRepository) *CompanyService {
	return &CompanyService{
		repository: repository,
	}
}

type RegisterCompanyRequest struct {
	Name string
}

type CompanyResponse struct {
	ID   uuid.UUID
	Name string
}

func (s *CompanyService) RegisterCompany(req RegisterCompanyRequest) (CompanyResponse, error) {
	company, _ := domain.RegisterCompany(req.Name)
	if err := s.repository.Save(company); err != nil {
		return CompanyResponse{}, err
	}
	return CompanyResponse{
		ID:   company.ID(),
		Name: company.Name(),
	}, nil
}

func (s *CompanyService) FetchCompanyByID(id uuid.UUID) (CompanyResponse, error) {
	company, err := s.repository.FindByID(id)
	if err != nil {
		return CompanyResponse{}, err
	}
	return CompanyResponse{
		ID:   company.ID(),
		Name: company.Name(),
	}, nil
}

type UpdateCompanyRequest struct {
	CompanyID uuid.UUID
	Name      string
}

func (s *CompanyService) UpdateCompany(req UpdateCompanyRequest) (CompanyResponse, error) {
	company, err := s.repository.FindByID(req.CompanyID)
	if err != nil {
		return CompanyResponse{}, err
	}

	company.Rename(req.Name)
	if err := s.repository.Save(&company); err != nil {
		return CompanyResponse{}, err
	}

	return CompanyResponse{
		ID:   company.ID(),
		Name: company.Name(),
	}, nil
}
