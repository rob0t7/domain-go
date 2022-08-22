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
	server := rest.New()
	url := "/companies/" + uuid.NewString()
	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateCompanySuccessfully(t *testing.T) {
	server := rest.New()
	companyName := "ACME INC"
	url := "/companies"
	reqBody := strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, companyName))

	r := httptest.NewRequest(http.MethodPost, url, reqBody)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, r)

	response := w.Result()
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, rest.JsonContentType, response.Header.Get("Content-Type"))
	var actualResponse rest.CompanyResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)
	assert.Equal(t, companyName, actualResponse.Name)
	assert.NotEqual(t, uuid.Nil, actualResponse.ID)
	id := actualResponse.ID
	location := response.Header.Get("Location")
	assert.Equal(t, fmt.Sprintf("/companies/%s", actualResponse.ID), location)

	r = httptest.NewRequest(http.MethodGet, location, nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, r)
	response = w.Result()

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, rest.JsonContentType, response.Header.Get("Content-Type"))
	err = json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)
	assert.Equal(t, companyName, actualResponse.Name)
	assert.Equal(t, id.String(), actualResponse.ID.String())

}

// func TestFetchCompanyCollection(t *testing.T) {
// 	url := "/companies"
// 	server := rest.New()
// 	r := httptest.NewRequest(http.MethodGet, url, nil)
// 	w := httptest.NewRecorder()
// 	server.ServeHTTP(w, r)
// 	assert.Equal(t, http.StatusOK, w.Code)
// }
