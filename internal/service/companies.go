package service

import (
	"context"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/google/uuid"
)

// CompanyRepository describes interface for company repository
type CompanyRepository interface {
	CreateCompany(ctx context.Context, company domain.Company) (uuid.UUID, error)
	GetCompany(ctx context.Context, id uuid.UUID) (domain.Company, error)
	SelectCompanies(ctx context.Context, limit, offset int) ([]domain.Company, error)
	DeleteCompany(ctx context.Context, id uuid.UUID) error
	UpdateCompany(ctx context.Context, uuid uuid.UUID, company domain.Company) error
}

// CompanyService is service to work with companies
type CompanyService struct {
	repo CompanyRepository
}

// NewCompanyService creates new company service
func NewCompanyService(repo CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

// Create creates new company
func (s *CompanyService) Create(ctx context.Context, company domain.Company) (uuid.UUID, error) {
	return s.repo.CreateCompany(ctx, company)
}

// Get gets company by id
func (s *CompanyService) Get(ctx context.Context, id uuid.UUID) (domain.Company, error) {
	return s.repo.GetCompany(ctx, id)
}

// Select selects list of companies
func (s *CompanyService) Select(ctx context.Context, limit, offset int) ([]domain.Company, error) {
	return s.repo.SelectCompanies(ctx, limit, offset)
}

// Delete deletes company by id
func (s *CompanyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCompany(ctx, id)
}

// Update updates company by id
func (s *CompanyService) Update(ctx context.Context, id uuid.UUID, company domain.Company) error {
	return s.repo.UpdateCompany(ctx, id, company)
}
