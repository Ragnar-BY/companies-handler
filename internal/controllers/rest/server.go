package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CompaniesUsecase interface {
	Create(ctx context.Context, company domain.Company) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Company, error)
	Select(ctx context.Context, limit, offset int) ([]domain.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, uuid uuid.UUID, company domain.Company) error
}

type UsersUsecase interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
}

type Server struct {
	srv *http.Server
	log *zap.Logger

	companies CompaniesUsecase
	users     UsersUsecase
}

// NewServer creates new server instance
func NewServer(addr string, log *zap.Logger, companies CompaniesUsecase, users UsersUsecase) *Server {
	e := gin.Default()

	srv := &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadHeaderTimeout: 60 * time.Second,
	}
	s := Server{
		srv:       srv,
		log:       log,
		companies: companies,
		users:     users,
	}
	s.routes(e)
	return &s
}

// Routes adds routes to server
func (s *Server) routes(e *gin.Engine) {
	e.GET("/healtz", s.Healtz)

	companies := e.Group("/companies")
	{
		companies.GET("/", s.SelectCompanies)
		companies.POST("/", s.CreateCompany)

		oneCompany := companies.Group("/:id")
		{
			oneCompany.GET("/", s.GetCompany)
			oneCompany.PATCH("/", s.UpdateCompany)
			oneCompany.DELETE("/", s.DeleteCompany)
		}
	}

	e.POST("/users", s.RegisterUser)
}

// Run starts server
func (s *Server) Run() error {
	err := s.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop stops server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Healtz(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
