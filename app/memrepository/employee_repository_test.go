package memrepository_test

import (
	"testing"

	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/memrepository"
)

func TestEmployeeRepository(t *testing.T) {
	companyRepo := memrepository.New()
	repo := memrepository.NewEmployeeRepository()
	app.RunEmployeeRepositoryTestSuite(t, repo, companyRepo)
}
