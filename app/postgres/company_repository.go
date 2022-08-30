package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/domain"
)

type CompanyRepository struct {
	db DBInterface
}

func NewCompanyRepository(db DBInterface) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) FindAll() []*domain.Company {
	var companies []*domain.Company
	rows, err := r.db.Query(`SELECT id,name FROM companies ORDER BY name`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id   uuid.UUID
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			panic(err)
		}
		companies = append(companies, domain.NewWithID(id, name))
	}
	rerr := rows.Close()
	if rerr != nil {
		panic(err)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return companies
}

func (r *CompanyRepository) FindByID(id uuid.UUID) (*domain.Company, error) {
	var (
		name string
	)
	err := r.db.QueryRow(`SELECT name FROM companies WHERE id = $1`, id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}
	return domain.NewWithID(id, name), nil
}

func (r *CompanyRepository) Insert(company *domain.Company) error {
	result, err := r.db.Exec(`INSERT INTO companies(id,name) VALUES( $1, $2)`, company.ID().String(), company.Name())
	if err != nil {
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
			return app.ErrConflict
		}
		return fmt.Errorf("failed to insert company record: %w", err)
	}
	rowsInserted, err := result.RowsAffected()
	if rowsInserted <= 0 {
		return fmt.Errorf("failed to insert company record")
	}
	if err != nil {
		return fmt.Errorf("failed to insert company record: %w", err)
	}
	return nil
}

func (r *CompanyRepository) Update(company *domain.Company) error {
	result, err := r.db.Exec(`UPDATE companies SET name = $1 WHERE id = $2`, company.Name(), company.ID().String())
	if err != nil {
		return err
	}
	rowsInserted, err := result.RowsAffected()
	if rowsInserted <= 0 {
		return app.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update company record: %w", err)
	}
	return nil
}

func (r *CompanyRepository) Delete(company *domain.Company) error {
	result, err := r.db.Exec(`DELETE FROM companies WHERE id = $1`, company.ID().String())
	if err != nil {
		return err
	}
	rowsInserted, err := result.RowsAffected()
	if rowsInserted <= 0 {
		return app.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to delete company record: %w", err)
	}
	return nil
}

func (r *CompanyRepository) Reset() {
	_, err := r.db.Exec(`TRUNCATE table companies CASCADE`)
	if err != nil {
		panic("failed to reset CompanyRepository")
	}
}
