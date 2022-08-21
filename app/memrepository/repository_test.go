package memrepository_test

import (
	"testing"

	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/memrepository"
)

func TestMemRepository(t *testing.T) {
	repo := memrepository.New()
	app.RunCompanyRepositoryTestSuite(t, repo)
}
