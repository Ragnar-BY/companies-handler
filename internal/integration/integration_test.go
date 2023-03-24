package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/Ragnar-BY/companies-handler/internal/broker"
	"github.com/Ragnar-BY/companies-handler/internal/config"
	"github.com/Ragnar-BY/companies-handler/internal/controllers/rest"
	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/Ragnar-BY/companies-handler/internal/repository/postgres"
	"github.com/Ragnar-BY/companies-handler/internal/service"
	"github.com/Ragnar-BY/companies-handler/internal/usecase"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type e2eTestSuite struct {
	suite.Suite

	srvAddr     string
	dbMigration *migrate.Migrate
	dbClient    *postgres.PostgresClient

	server *rest.Server
}

func Test_E2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	logger, _ := zap.NewDevelopment()

	cfg, err := config.LoadConfig(".env")
	s.Require().NoError(err)
	s.srvAddr = cfg.ServerAddress

	pgSettings := postgres.PostgresSettings{
		Addr:     cfg.PostgresAddress,
		Username: cfg.PostgresUser,
		Password: cfg.PostgresPassword,
		Database: cfg.PostgresDB,
	}

	dbConn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgSettings.Username, pgSettings.Password, pgSettings.Addr, pgSettings.Database)
	dbClient, err := postgres.NewPostgresClient(pgSettings)
	s.Require().NoError(err)
	s.dbClient = dbClient

	msgBroker := broker.NewBroker()
	eventSrv := service.NewEventSender(msgBroker)

	companySrv := service.NewCompanyService(dbClient)
	companyUsecase := usecase.NewCompanyUsecase(companySrv, eventSrv)

	userSrv := service.NewUserService(dbClient)
	authSrv := service.NewAuthService([]byte(cfg.JWTKey))
	authUsecase := usecase.NewAuthUsecase(authSrv, userSrv)
	srv := rest.NewServer(cfg.ServerAddress, logger, companyUsecase, authUsecase)
	s.server = srv

	s.dbMigration, err = migrate.New("file://../../postgres/migrations", dbConn)
	s.Require().NoError(err)
	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}

	go srv.Run()
}

func (s *e2eTestSuite) TearDownSuite() {
	p, err := os.FindProcess(syscall.Getpid())
	s.Require().NoError(err)
	err = p.Signal(syscall.SIGINT)
	s.Require().NoError(err)
}

func (s *e2eTestSuite) SetupTest() {
	if err := s.dbMigration.Up(); err != nil && err != migrate.ErrNoChange {
		s.Require().NoError(err)
	}
}

func (s *e2eTestSuite) TearDownTest() {
	s.NoError(s.dbMigration.Down())
}

func (s *e2eTestSuite) Test_EndToEnd_GetCompany() {
	company := domain.Company{
		Name:              "test-company",
		Description:       "some description",
		AmountOfEmployees: 123,
		Registered:        true,
		Type:              domain.Corporations,
	}
	ctx := context.Background()

	id, err2 := s.dbClient.CreateCompany(ctx, company)
	s.Require().NoError(err2)

	testCases := []struct {
		name         string
		id           uuid.UUID
		statusCode   int
		responseBody string
	}{
		{
			name:         "valid test",
			id:           id,
			statusCode:   http.StatusOK,
			responseBody: fmt.Sprintf(`{"id":"%v","name":"test-company","description":"some description","amount_of_employees":123,"registered":true,"type":"Corporations"}`, id),
		},
		{
			name:         "wrong id",
			id:           uuid.Nil,
			statusCode:   http.StatusNotFound,
			responseBody: `{"error":"can not get company: sql: no rows in result set"}`,
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/companies/%v", s.srvAddr, tt.id), http.NoBody)
			s.NoError(err)

			req.Header.Set("Content-Type", "application/json")

			client := http.Client{}
			response, err := client.Do(req)
			s.NoError(err)
			s.Equal(tt.statusCode, response.StatusCode)

			byteBody, err := io.ReadAll(response.Body)
			s.NoError(err)

			s.Equal(tt.responseBody, string(byteBody))
			response.Body.Close()
		})
	}

	err := s.dbClient.DeleteCompany(ctx, id)
	s.NoError(err)
}

func (s *e2eTestSuite) Test_EndToEnd_SelectCompanies() {
	company := domain.Company{
		ID:                uuid.New(),
		Name:              "test-company",
		Description:       "some description",
		AmountOfEmployees: 123,
		Registered:        true,
		Type:              domain.Corporations,
	}
	ctx := context.Background()

	id, err := s.dbClient.CreateCompany(ctx, company)
	s.NoError(err)

	testCases := []struct {
		name         string
		offset       int
		statusCode   int
		responseBody string
	}{
		{
			name:         "valid test",
			offset:       0,
			statusCode:   http.StatusOK,
			responseBody: fmt.Sprintf(`[{"id":"%v","name":"test-company","description":"some description","amount_of_employees":123,"registered":true,"type":"Corporations"}]`, id),
		},
		{
			name:         "null output",
			offset:       3,
			statusCode:   http.StatusOK,
			responseBody: "[]",
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/companies?offset=%d", s.srvAddr, tt.offset), http.NoBody)
			s.NoError(err)

			req.Header.Set("Content-Type", "application/json")

			client := http.Client{}
			response, err := client.Do(req)
			s.NoError(err)
			s.Equal(tt.statusCode, response.StatusCode)

			byteBody, err := io.ReadAll(response.Body)
			s.NoError(err)

			s.Equal(tt.responseBody, string(byteBody))
			response.Body.Close()
		})
	}
}

func (s *e2eTestSuite) Test_EndToEnd_CreateCompany() {
	company := domain.Company{
		Name:              "company-create",
		Description:       "some description",
		AmountOfEmployees: 123,
		Registered:        true,
		Type:              domain.Corporations,
	}
	ctx := context.Background()

	userBody := `{
		"username":"test user",
		"email": "test@test.com",
		"password": "12345678"
	}`
	UserReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/register", s.srvAddr), strings.NewReader(userBody))
	s.NoError(err)
	UserReq.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	userResponse, err := client.Do(UserReq)
	s.NoError(err)
	s.Equal(http.StatusCreated, userResponse.StatusCode)

	token := struct {
		Token string
	}{}
	err = json.NewDecoder(userResponse.Body).Decode(&token)
	s.NoError(err)
	userResponse.Body.Close()

	testCases := []struct {
		name          string
		statusCode    int
		token         string
		companyCreate bool
		companyBody   string
	}{
		{
			name:       "valid test",
			statusCode: http.StatusCreated,
			token:      token.Token,
			companyBody: `{
				"name": "company-create",
				"description": "some description",
				"amount_of_employees":123,
				"registered": true,
				"type": "Corporations"
			}`,
			companyCreate: true,
		},
		{
			name:          "unauthorized",
			statusCode:    http.StatusUnauthorized,
			token:         "",
			companyCreate: false,
			companyBody:   "",
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/companies", s.srvAddr), strings.NewReader(tt.companyBody))
			s.NoError(err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			response, err := client.Do(req)
			s.NoError(err)

			s.Require().Equal(tt.statusCode, response.StatusCode)

			if tt.companyCreate {
				idStr := struct {
					ID uuid.UUID
				}{}
				err = json.NewDecoder(response.Body).Decode(&idStr)
				s.NoError(err)

				fmt.Println(idStr.ID)
				company.ID = idStr.ID
				cmp, err := s.dbClient.GetCompany(ctx, idStr.ID)
				s.NoError(err)
				s.Equal(cmp, company)

				err = s.dbClient.DeleteCompany(ctx, idStr.ID)
				s.NoError(err)
			}

			response.Body.Close()
		})
	}
}

func (s *e2eTestSuite) Test_EndToEnd_DeleteCompany() {
	ctx := context.Background()

	company := domain.Company{
		Name:              "test-company",
		Description:       "some description",
		AmountOfEmployees: 123,
		Registered:        true,
		Type:              domain.Corporations,
	}

	id, err2 := s.dbClient.CreateCompany(ctx, company)
	s.Require().NoError(err2)

	userBody := `{
		"username":"test user",
		"email": "test@test.com",
		"password": "12345678"
	}`
	UserReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/register", s.srvAddr), strings.NewReader(userBody))
	s.NoError(err)
	UserReq.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	userResponse, err := client.Do(UserReq)
	s.NoError(err)
	s.Equal(http.StatusCreated, userResponse.StatusCode)

	token := struct {
		Token string
	}{}
	err = json.NewDecoder(userResponse.Body).Decode(&token)
	s.NoError(err)
	userResponse.Body.Close()

	testCases := []struct {
		name       string
		statusCode int
		token      string
		id         uuid.UUID
	}{
		{
			name:       "valid test",
			statusCode: http.StatusOK,
			token:      token.Token,
			id:         id,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			token:      "",
		},
	}

	for _, tt := range testCases {
		s.Run(tt.name, func() {
			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("http://%s/companies/%v", s.srvAddr, tt.id), http.NoBody)
			s.NoError(err)

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))

			response, err := client.Do(req)
			s.NoError(err)

			s.Require().Equal(tt.statusCode, response.StatusCode)
			if tt.statusCode == http.StatusOK {
				response.Body.Close()

				_, err = s.dbClient.GetCompany(ctx, id)
				s.Equal("can not get company: sql: no rows in result set", err.Error())
			}
		})
	}
}
