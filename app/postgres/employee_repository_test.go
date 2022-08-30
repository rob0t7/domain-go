package postgres_test

import (
	"testing"

	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/postgres"
	"github.com/stretchr/testify/require"
)

func TestEmployeeRepository(t *testing.T) {
	db, err := postgres.Open("postgres://postgres:postgres@localhost:5432/app")
	require.NoError(t, err)
	defer db.Close() // nolint:staticcheck // don't need to verify db was closed in test
	repo := postgres.NewEmployeeRepository(db)
	companyRepo := postgres.NewCompanyRepository(db)
	app.RunEmployeeRepositoryTestSuite(t, repo, companyRepo)
}
