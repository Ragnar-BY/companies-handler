package rest

import (
	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/google/uuid"
)

type company struct {
	ID                uuid.UUID
	Name              string
	Description       string
	AmountOfEmployees int
	Registered        bool
	Type              string
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
