package postgres

import (
	"database/sql"
	"errors"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	postgreslib "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/domain"
)

type PostgresDB struct {
	*sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	var pg PostgresDB
	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/app?sslmode=disable")
	if err != nil {
		return nil, err
	}
	pg.DB = db
	driver, err := postgreslib.WithInstance(pg.DB, &postgreslib.Config{})
	if err != nil {
		return nil, err
	}

	_, cwd, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(cwd), "..", "..", "migrations")
	migrator, err := migrate.NewWithDatabaseInstance("file://"+path, "postgres", driver)
	if err != nil {
		return nil, err
	}
	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return &pg, nil
}

type PosgtresCompanyRepository struct {
	db *PostgresDB
}

func NewCompanyRepository(db *PostgresDB) *PosgtresCompanyRepository {
	return &PosgtresCompanyRepository{
		db: db,
	}
}

func (r *PosgtresCompanyRepository) FindAll() []*domain.Company {
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
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return companies
}

func (r *PosgtresCompanyRepository) FindByID(id uuid.UUID) (*domain.Company, error) {
	var (
		name string
	)
	sqlQuery := `SELECT name FROM companies WHERE id = $1`
	err := r.db.QueryRow(sqlQuery, id.String()).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}
	return domain.NewWithID(id, name), nil
}

func (r *PosgtresCompanyRepository) Insert(company *domain.Company) error {
	sql := `INSERT INTO companies(id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(sql, company.ID().String(), company.Name())
	if err != nil {
		return err
	}
	return nil
}

func (r *PosgtresCompanyRepository) Update(company *domain.Company) error {
	result, err := r.db.Exec(`UPDATE companies SET name = $1 WHERE id = $2`, company.Name(), company.ID().String())
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return app.ErrNotFound
	}
	return nil
}
func (r *PosgtresCompanyRepository) Delete(company *domain.Company) error {
	sql := `DELETE FROM companies WHERE id = $1`
	result, err := r.db.Exec(sql, company.ID().String())
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return app.ErrNotFound
	}
	return nil
}

func (r *PosgtresCompanyRepository) Reset() {
	sql := `TRUNCATE TABLE companies CASCADE`
	_, err := r.db.Exec(sql)
	if err != nil {
		panic(err)
	}
}
