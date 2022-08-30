package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rob0t7/domain-go/app"
	"github.com/rob0t7/domain-go/app/domain"
)

type DBInterface interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
type EmployeeRepository struct {
	db DBInterface
}

func NewEmployeeRepository(db DBInterface) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (repo *EmployeeRepository) FindAll(companyID uuid.UUID) []*domain.Employee {
	var employees []*domain.Employee
	rows, err := repo.db.Query(`SELECT id,first_name,last_name,salary FROM employees WHERE company_id = $1 ORDER BY last_name,first_name`, companyID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        uuid.UUID
			firstName string
			lastName  string
			salary    int64
		)
		if err := rows.Scan(&id, &firstName, &lastName, &salary); err != nil {
			panic(err)
		}
		employees = append(employees, domain.NewEmployee(id, companyID, firstName, lastName, salary))
	}
	if rerr := rows.Close(); rerr != nil {
		panic(rerr)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return employees
}

func (repo *EmployeeRepository) FindByID(id uuid.UUID) (*domain.Employee, error) {
	var (
		firstName string
		lastName  string
		companyID uuid.UUID
		salary    int64
	)
	err := repo.db.QueryRow(
		`SELECT first_name,last_name,company_id,salary FROM employees WHERE id = $1`,
		id,
	).Scan(&firstName, &lastName, &companyID, &salary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}
	return domain.NewEmployee(id, companyID, firstName, lastName, salary), nil
}

func (repo *EmployeeRepository) Reset() {
	if _, err := repo.db.Exec(`TRUNCATE TABLE employees`); err != nil {
		panic("failed to reset EmployeRepository")
	}
}

func (repo *EmployeeRepository) Insert(employee *domain.Employee) error {
	result, err := repo.db.Exec(
		`INSERT INTO employees(id, first_name, last_name, company_id, salary) VALUES($1, $2, $3, $4, $5)`,
		employee.ID(),
		employee.FirstName(),
		employee.LastName(),
		employee.CompanyID(),
		0,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected <= 0 {
		return fmt.Errorf("failed to insert employee")
	}
	return nil
}

func (repo *EmployeeRepository) Delete(employee *domain.Employee) error {
	results, err := repo.db.Exec(`DELETE FROM employees WHERE id = $1`, employee.ID())
	if err != nil {
		return err
	}
	count, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return app.ErrNotFound
	}
	return nil
}
