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
	if err := s.repository.Insert(company); err != nil {
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
	if err := s.repository.Update(company); err != nil {
		return CompanyResponse{}, err
	}

	return CompanyResponse{
		ID:   company.ID(),
		Name: company.Name(),
	}, nil
}

type DeleteCompanyRequest struct {
	CompanyID uuid.UUID
}

func (s *CompanyService) DeleteCompany(req DeleteCompanyRequest) error {
	company, err := s.repository.FindByID(req.CompanyID)
	if err != nil {
		return err
	}
	return s.repository.Delete(company)
}

type CompanyCollectionResponse struct {
	Total     int
	Companies []CompanyResponse
}

func (s *CompanyService) FetchAll() (response CompanyCollectionResponse) {
	companies := s.repository.FindAll()
	response.Total = len(companies)
	for _, company := range companies {
		response.Companies = append(response.Companies, CompanyResponse{company.ID(), company.Name()})
	}
	return response
}
