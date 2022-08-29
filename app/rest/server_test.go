package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/memrepository"
	"github.com/rob0t7/domain-go/app/rest"
	"github.com/stretchr/testify/assert"
)

func newRESTServer() *rest.RESTServer {
	return rest.New(app.NewCompanyService(memrepository.New()))
}
func TestNotFound(t *testing.T) {
	server := newRESTServer()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	server.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
