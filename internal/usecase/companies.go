package usecase

import (
	"context"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/google/uuid"
)

const (
	createCompanyTopic = "create-company"
	updateCompanyTopic = "update-company"
	deleteCompanyTopic = "delete-company"
)

// CompanyService describes interface for company service
type CompanyService interface {
	Create(ctx context.Context, company domain.Company) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Company, error)
	Select(ctx context.Context, limit, offset int) ([]domain.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, uuid uuid.UUID, company domain.Company) error
}

// EventService describe service for sending events in message broker
type EventService interface {
	SendEvent(topic string, message any) error
}

// CompanyUsecase is usecase for companies
type CompanyUsecase struct {
	srv    CompanyService
	events EventService
}

// NewCompanyUsecase creates new company usecase
func NewCompanyUsecase(companies CompanyService, events EventService) *CompanyUsecase {
	return &CompanyUsecase{srv: companies, events: events}
}

// Create creates new company
func (u *CompanyUsecase) Create(ctx context.Context, company domain.Company) (uuid.UUID, error) {
	id, err := u.srv.Create(ctx, company)
	if err != nil {
		return uuid.Nil, err
	}
	err = u.events.SendEvent(createCompanyTopic, id)
	return id, err
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
	err := u.srv.Delete(ctx, id)
	if err != nil {
		return err
	}
	return u.events.SendEvent(deleteCompanyTopic, id)
}

// Update updates company
func (u *CompanyUsecase) Update(ctx context.Context, id uuid.UUID, company domain.Company) error {
	company.ID = id
	err := u.srv.Update(ctx, id, company)

	if err != nil {
		return err
	}
	return u.events.SendEvent(updateCompanyTopic, company)
}
