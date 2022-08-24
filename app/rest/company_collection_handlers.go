package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rob0t7/domain-go/app"
)

func (s *RESTServer) CompanyCollectionHandler(w http.ResponseWriter, r *http.Request) {
	if !listCompanyRE.MatchString(r.URL.Path) {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodPost {
		s.CreateCompanyHandler(w, r)
		return
	}
	if r.Method == http.MethodGet {
		s.FetchCompanyCollection(w, r)
		return
	}
	w.Header().Set("Allow", "GET POST")
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *RESTServer) CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var registerCompanyReq app.RegisterCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&registerCompanyReq); err != nil {
		panic(err) // TODO: Proper error handling
	}

	createResponse, _ := s.companyService.RegisterCompany(registerCompanyReq)
	response := CompanyResponse{
		ID:   createResponse.ID,
		Name: createResponse.Name,
	}
	w.Header().Set("Content-Type", JsonContentType)
	w.Header().Set("Location", fmt.Sprintf(`/companies/%s`, response.ID))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		panic(err) // TODO: Proper error handling
	}
}

type CompanyCollectionResponse struct {
	Total     int               `json:"total"`
	Companies []CompanyResponse `json:"companies"`
}

func (s *RESTServer) FetchCompanyCollection(w http.ResponseWriter, r *http.Request) {
	var response CompanyCollectionResponse
	response.Companies = make([]CompanyResponse, 0)

	companies := s.companyService.FetchAll()
	response.Total = companies.Total
	for _, company := range companies.Companies {
		response.Companies = append(response.Companies, CompanyResponse{ID: company.ID, Name: company.Name})
	}

	w.Header().Set("Content-Type", JsonContentType)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		panic(err)
	}
}
