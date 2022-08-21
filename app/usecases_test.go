package app_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/memrepository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterCompany(t *testing.T) {
	companyService := app.NewCompanyService(memrepository.New())

	req := app.RegisterCompanyRequest{Name: "ACME INC"}
	res, err := companyService.RegisterCompany(req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, res.ID)
	assert.Equal(t, "ACME INC", res.Name)

	company, err := companyService.FetchCompanyByID(res.ID)
	require.NoError(t, err)
	assert.Equal(t, res, company)
}

func TestFetchCompany(t *testing.T) {
	companyService := app.NewCompanyService(memrepository.New())
	company, _ := companyService.RegisterCompany(app.RegisterCompanyRequest{Name: "ACME INC"})

	t.Run("Successfully fetch company", func(t *testing.T) {
		res, err := companyService.FetchCompanyByID(company.ID)
		require.NoError(t, err)
		assert.Equal(t, company, res)

	})

	t.Run("return ErrNotFound if company does not exist", func(t *testing.T) {
		_, err := companyService.FetchCompanyByID(uuid.New())
		require.ErrorIs(t, err, app.ErrNotFound)
	})
}

func TestUpdateCompanyName(t *testing.T) {
	companyService := app.NewCompanyService(memrepository.New())
	company, err := companyService.RegisterCompany(app.RegisterCompanyRequest{Name: "ACME INC"})
	require.NoError(t, err)

	t.Run("returns ErrNotFound when the company does not exist", func(t *testing.T) {
		_, err := companyService.UpdateCompany(app.UpdateCompanyRequest{CompanyID: uuid.New()})
		require.ErrorIs(t, err, app.ErrNotFound)
	})

	t.Run("updates the company name accordingly", func(t *testing.T) {
		newName := "Wiley Industries"
		res, err := companyService.UpdateCompany(
			app.UpdateCompanyRequest{
				CompanyID: company.ID,
				Name:      newName,
			},
		)
		require.NoError(t, err)
		require.Equal(t, newName, res.Name)
	})

}
