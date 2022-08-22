package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rob0t7/domain-go/app/rest"
	"github.com/stretchr/testify/assert"
)

func TestNotFound(t *testing.T) {
	server := rest.New()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	server.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
