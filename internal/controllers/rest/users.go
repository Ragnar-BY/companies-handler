package rest

import (
	"net/http"

	"github.com/Ragnar-BY/companies-handler/internal/domain"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type user struct {
	Username string
	Email    string
	Password string
}

// RegisterUser saves user with hash password
func (s *Server) RegisterUser(c *gin.Context) {
	var u user
	err := c.ShouldBind(&u)
	if err != nil {
		s.log.Error("can not bind user", zap.Error(err))
		c.JSON(http.StatusBadRequest, err)
	}

	newUser := domain.User{
		Username: u.Username,
		Email:    u.Email,
	}
	err = newUser.HashPassword(u.Password)
	if err != nil {
		s.log.Error("can not hash user password", zap.Error(err))
		c.JSON(http.StatusBadRequest, err)
	}

	us, err := s.users.CreateUser(c.Request.Context(), newUser)
	if err != nil {
		s.log.Error("can not create new user", zap.Error(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, us)
}
