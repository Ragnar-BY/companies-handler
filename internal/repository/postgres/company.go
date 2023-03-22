package postgres

import (
	"context"
	"time"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/google/uuid"
)

type company struct {
	ID                uuid.UUID `db:"id"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	AmountOfEmployees int       `db:"amount_of_employees"`
	Registered        bool      `db:"registered"`
	Type              string    `db:"type"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func fromDomain(c domain.Company) company {
	return company{
		ID:                c.ID,
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		Type:              string(c.Type),
	}
}

func toDomain(c company) domain.Company {
	return domain.Company{
		ID:                c.ID,
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		Type:              domain.CompanyType(c.Type),
	}
}

// Create created new company in database
func (c *PostgresClient) Create(ctx context.Context, company domain.Company) (uuid.UUID,error) {
	cmp := fromDomain(company)
	var id uuid.UUID
	stmt,err:=c.db.PrepareNamedContext(ctx, `INSERT INTO companies( name,description, amount_of_employees, registered,type) 
	VALUES (:name, :description, :amount_of_employees, :registered, :type) RETURNING id `)
	if err!= nil{
		return uuid.Nil,err
	}
	err = stmt.GetContext(ctx,&id,cmp)
	return id,err	
}

// Get gets company from DB by id
func (c *PostgresClient) Get(ctx context.Context, id uuid.UUID) (domain.Company, error) {
	var cmp company
	err := c.db.GetContext(ctx, &cmp, "SELECT * FROM companies WHERE id=$1", id)
	if err != nil {
		return domain.Company{}, err
	}
	return toDomain(cmp), nil
}

// Select selects all companies from database
func (c *PostgresClient) Select(ctx context.Context, limit,offset int) ([]domain.Company, error) {
	var cmps []company
	err := c.db.SelectContext(ctx,&cmps, "SELECT * FROM companies LIMIT $1 OFFSET $2",limit,offset)
	if err != nil {
		return nil, err
	}
	companies := make([]domain.Company, 0)
	for _, cmp := range cmps {
		companies = append(companies, toDomain(cmp))
	}
	return companies, nil
}

// Delete deletes company from DB by id
func (c *PostgresClient) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := c.db.ExecContext(ctx, "DELETE FROM companies WHERE id=$1", id)
	return err
}

// Update updates company by id
func (c *PostgresClient) Update(ctx context.Context, uuid uuid.UUID, company domain.Company) error {
	cmp := fromDomain(company)
	cmp.ID = uuid
	cmp.UpdatedAt = time.Now()
	_, err := c.db.NamedExecContext(ctx,
		`UPDATE companies SET name=:name, description=:description, amount_of_employees=:amount_of_employees, 
		registered=:registered,type=:type, updated_at=:updated_at
		 WHERE id=:id`, cmp)
	return err
}
