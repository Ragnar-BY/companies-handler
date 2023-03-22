package domain

import "github.com/google/uuid"


type CompanyType string

const (
	Corporations CompanyType = "Corporations"
	NonProfit CompanyType = "NonProfit"
	Cooperative CompanyType = "Cooperative"
	SoleProprietorship CompanyType = "Sole Proprietorship"
)

// Company is company description
type Company struct {
	ID uuid.UUID
	Name string
	Description string
	AmountOfEmployees int
	Registered bool
	Type CompanyType
}