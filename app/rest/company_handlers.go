package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
)

const JsonContentType = "application/hal+json"

var listCompanyRE = regexp.MustCompile(`^/companies$`)
var getCompanyRE = regexp.MustCompile(`^\/companies\/([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$`)

type CompanyResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (s *RESTServer) CompanyInstanceHandler(w http.ResponseWriter, r *http.Request) {
	if getCompanyRE.MatchString(r.URL.Path) {
		matches := getCompanyRE.FindStringSubmatch(r.URL.Path)
		id := uuid.MustParse(matches[1])

		switch r.Method {
		case http.MethodGet:
			s.FetchCompany(w, r, id)
		case http.MethodDelete:
			s.DeleteCompany(w, r, id)
		case http.MethodPut:
			s.UpdateCompany(w, r, id)
		default:
			w.Header().Set("Allow", "GET PUT DELETE")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	http.NotFound(w, r)
}

func (s *RESTServer) FetchCompany(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	company, err := s.companyService.FetchCompanyByID(id)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := CompanyResponse{
		ID:   company.ID,
		Name: company.Name,
	}
	w.Header().Set("Content-Type", JsonContentType)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func (s *RESTServer) DeleteCompany(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	err := s.companyService.DeleteCompany(app.DeleteCompanyRequest{CompanyID: id})
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *RESTServer) UpdateCompany(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var updateRequest CompanyResponse
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed json unmarshalling: %v", err), http.StatusBadRequest)
		return
	}
	company, err := s.companyService.UpdateCompany(app.UpdateCompanyRequest{CompanyID: id, Name: updateRequest.Name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := CompanyResponse{
		ID:   company.ID,
		Name: company.Name,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
