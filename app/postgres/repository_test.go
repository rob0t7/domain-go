package postgres_test

import (
	"testing"

	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/postgres"
	"github.com/stretchr/testify/require"
)

func TestPostgresCompanyRepository(t *testing.T) {
	db, err := postgres.NewPostgresDB()
	require.NoError(t, err)
	repo := postgres.NewCompanyRepository(db)
	app.RunCompanyRepositoryTestSuite(t, repo)
}
