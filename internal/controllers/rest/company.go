package rest

import (
	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/google/uuid"
)

type company struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name" validate:"required,min=4,max=15"`
	Description       string    `json:"description" validate:"max=3000"`
	AmountOfEmployees int       `json:"amount_of_employees" validate:"required,min=1"`
	Registered        bool      `json:"registered" validate:"required"`
	Type              string    `json:"type" validate:"required,oneof='Corporations' 'NonProfit' 'Cooperative' 'Sole Proprietorship'"`
}

func companyToDomain(c company) domain.Company {
	return domain.Company{
		ID:                c.ID,
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		Type:              domain.CompanyType(c.Type),
	}
}

func domainToCompany(c domain.Company) company {
	return company{
		ID:                c.ID,
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		Type:              string(c.Type),
	}
}
