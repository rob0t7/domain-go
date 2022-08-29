package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrieveReturns404IfNotFound(t *testing.T) {
	server := newRESTServer()
	url := "/companies/" + uuid.NewString()
	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateCompanySuccessfully(t *testing.T) {
	server := newRESTServer()
	companyName := "ACME INC"
	createResponse, location := createCompany(t, server, companyName)
	assert.Equal(t, companyName, createResponse.Name)
	assert.NotEqual(t, uuid.Nil, createResponse.ID)

	r := httptest.NewRequest(http.MethodGet, location, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, rest.JsonContentType, response.Header.Get("Content-Type"))
	var actualResponse rest.CompanyResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)
	assert.Equal(t, companyName, actualResponse.Name)
	assert.Equal(t, createResponse.ID, actualResponse.ID)
}

func createCompany(t *testing.T, server *rest.RESTServer, name string) (rest.CompanyResponse, string) {
	t.Helper()
	reqBody := strings.NewReader(fmt.Sprintf(`{"name": "%s"}`, name))
	r := httptest.NewRequest(http.MethodPost, "/companies", reqBody)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, r)

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, rest.JsonContentType, w.Header().Get("Content-Type"))
	var actualResponse rest.CompanyResponse
	err := json.NewDecoder(w.Result().Body).Decode(&actualResponse)
	require.NoError(t, err)
	location := w.Header().Get("Location")
	assert.Equal(t, fmt.Sprintf("/companies/%s", actualResponse.ID), location)

	return actualResponse, location
}
func TestFetchCompanyCollection(t *testing.T) {
	url := "/companies"
	server := newRESTServer()
	company1, _ := createCompany(t, server, "Wiley Coyote LTD")
	company2, _ := createCompany(t, server, "ACME INC")
	expectedResp := rest.CompanyCollectionResponse{
		Total:     2,
		Companies: []rest.CompanyResponse{company2, company1},
	}

	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/hal+json", w.Header().Get("Content-Type"))
	var actualResponse rest.CompanyCollectionResponse
	err := json.NewDecoder(w.Result().Body).Decode(&actualResponse)
	require.NoError(t, err)
	assert.Equal(t, expectedResp, actualResponse)
}

func TestDeleteCompany(t *testing.T) {
	server := newRESTServer()
	company, _ := createCompany(t, server, "ACME INC")

	r := httptest.NewRequest(http.MethodDelete, "/companies/"+company.ID.String(), nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = httptest.NewRecorder()
	server.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateCompany(t *testing.T) {
	server := newRESTServer()
	company, _ := createCompany(t, server, "ACME INC")

	url := "/companies/" + company.ID.String()
	reqBody := strings.NewReader(`{"name": "Wiley Coyote LTD"}`)
	r := httptest.NewRequest(http.MethodPut, url, reqBody)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Allow", "application/json")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualResponse rest.CompanyResponse
	err := json.NewDecoder(w.Result().Body).Decode(&actualResponse)
	require.NoError(t, err)
	assert.Equal(t, company.ID, actualResponse.ID)
	assert.Equal(t, "Wiley Coyote LTD", actualResponse.Name)

}
