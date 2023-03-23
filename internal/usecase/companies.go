package usecase

import (
	"context"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/google/uuid"
)

// CompanyService describes interface for company service
type CompanyService interface {
	Create(ctx context.Context, company domain.Company) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Company, error)
	Select(ctx context.Context, limit, offset int) ([]domain.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, uuid uuid.UUID, company domain.Company) error
}

// CompanyUsecase is usecase for companies
type CompanyUsecase struct {
	srv CompanyService
}

// NewCompanyUsecase creates new company usecase
func NewCompanyUsecase(companies CompanyService) *CompanyUsecase {
	return &CompanyUsecase{srv: companies}
}

// Create creates new company
func (u *CompanyUsecase) Create(ctx context.Context, company domain.Company) (uuid.UUID, error) {
	return u.srv.Create(ctx, company)
}

// Get gets new company
func (u *CompanyUsecase) Get(ctx context.Context, id uuid.UUID) (domain.Company, error) {
	return u.srv.Get(ctx, id)
}

// Select selects list of companies
func (u *CompanyUsecase) Select(ctx context.Context, limit, offset int) ([]domain.Company, error) {
	return u.srv.Select(ctx, limit, offset)
}

// Delete deletes company
func (u *CompanyUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.srv.Delete(ctx, id)
}

// Update updates company
func (u *CompanyUsecase) Update(ctx context.Context, id uuid.UUID, company domain.Company) error {
	return u.srv.Update(ctx, id, company)
}
