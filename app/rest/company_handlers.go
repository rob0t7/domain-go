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

func (s *RESTServer) CompanyCollectionHandler(w http.ResponseWriter, r *http.Request) {
	if !listCompanyRE.MatchString(r.URL.Path) {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodPost {
		s.CreateCompanyHandler(w, r)
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
	json.NewEncoder(w).Encode(&response)
}

func (s *RESTServer) CompanyInstanceHandler(w http.ResponseWriter, r *http.Request) {
	if getCompanyRE.MatchString(r.URL.Path) {
		matches := getCompanyRE.FindStringSubmatch(r.URL.Path)
		id := uuid.MustParse(matches[1])

		switch r.Method {
		case http.MethodGet:
			s.FetchCompany(w, r, id)
		default:
			w.Header().Set("Allow", "GET")
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
	json.NewEncoder(w).Encode(&response)
}
