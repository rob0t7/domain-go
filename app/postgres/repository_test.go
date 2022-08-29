package postgres_test

import (
	"testing"

	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/postgres"
	"github.com/stretchr/testify/require"
)

func TestPostgresCompanyRepository(t *testing.T) {
	db, err := postgres.Open("postgres://postgres:postgres@localhost:5432/app")
	defer db.Close() // nolint:staticcheck // don't need to verify db was closed in test
	require.NoError(t, err)
	repo := postgres.NewCompanyRepository(db)
	app.RunCompanyRepositoryTestSuite(t, repo)
}
