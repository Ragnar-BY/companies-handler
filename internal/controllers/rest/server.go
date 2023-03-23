package rest

import (
	"context"
	"net/http"
	"strings"
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

type AuthUsecase interface {
	SignUp(ctx context.Context, user domain.User) (string, error)
	SignIn(ctx context.Context, email string, password string) (string, error)
	ValidateToken(signedToken string) error
}

type Server struct {
	srv *http.Server
	log *zap.Logger

	companies CompaniesUsecase
	auth      AuthUsecase
}

// NewServer creates new server instance
func NewServer(addr string, log *zap.Logger, companies CompaniesUsecase, auth AuthUsecase) *Server {
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
		auth:      auth,
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
		companies.POST("/", s.CreateCompany).Use(s.Auth())

		oneCompany := companies.Group("/:id")
		{
			oneCompany.GET("/", s.GetCompany)
			oneCompany.PATCH("/", s.UpdateCompany).Use(s.Auth())
			oneCompany.DELETE("/", s.DeleteCompany).Use(s.Auth())
		}
	}

	e.POST("/register", s.RegisterUser)
	e.POST("/signin", s.SignIn)
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

// Auth is middleware to auth
func (s *Server) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is not bearer"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		err := s.auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
