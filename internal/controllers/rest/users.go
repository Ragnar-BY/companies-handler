package rest

import (
	"net/http"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type user struct {
	Username string `validate:"required,min=4"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

// RegisterUser saves user with hash password
func (s *Server) RegisterUser(c *gin.Context) {
	var u user
	err := c.ShouldBind(&u)
	if err != nil {
		s.log.Error("can not bind user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = validate.Struct(&u)
	if err != nil {
		s.log.Error("can not validate user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser := domain.User{
		Username: u.Username,
		Email:    u.Email,
	}
	err = newUser.HashPassword(u.Password)
	if err != nil {
		s.log.Error("can not hash user password", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.auth.SignUp(c.Request.Context(), newUser)
	if err != nil {
		s.log.Error("can not create new user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// SignIn sign in user
func (s *Server) SignIn(c *gin.Context) {
	var u user
	err := c.ShouldBind(&u)
	if err != nil {
		s.log.Error("can not bind user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.auth.SignIn(c.Request.Context(), u.Email, u.Password)
	if err != nil {
		s.log.Error("can not sign in", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
