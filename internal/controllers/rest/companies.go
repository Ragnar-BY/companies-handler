package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	defaultLimit = 20
)

// GetCompany gets company by id from path param
func (s *Server) GetCompany(c *gin.Context) {
	paramID := c.Param("id")
	id, err := uuid.Parse(paramID)
	if err != nil {
		s.log.Error("can not parse id", zap.Error(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}
	company, err := s.companies.Get(c.Request.Context(), id)
	if err != nil {
		s.log.Error("can not get company", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, company)
}

// SelectCompanies selects list of company according query params limit and offset
func (s *Server) SelectCompanies(c *gin.Context) {
	limitQuery, offsetQuery := c.Query("limit"), c.Query("offset")
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		s.log.Error("can not parse limit", zap.Error(err))
		limit = defaultLimit
	}
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		s.log.Error("can not parse offset", zap.Error(err))
		offset = 0
	}
	companies, err := s.companies.Select(c.Request.Context(), limit, offset)
	if err != nil {
		s.log.Error("can not select companies", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, companies)
}

// DeleteCompany deletes company by id from path param
func (s *Server) DeleteCompany(c *gin.Context) {
	paramID := c.Param("id")
	id, err := uuid.Parse(paramID)
	if err != nil {
		s.log.Error("can not parse id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = s.companies.Delete(c.Request.Context(), id)
	if err != nil {
		s.log.Error("can not delete company", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "")
}

// CreateCompany creates new company
func (s *Server) CreateCompany(c *gin.Context) {
	var cmp company
	err := c.ShouldBind(&cmp)
	if err != nil {
		s.log.Error("can not bind company", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	newCompany := companyToDomain(cmp)
	id, err := s.companies.Create(c.Request.Context(), newCompany)
	if err != nil {
		s.log.Error("can not create new company", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// Updatecompany updates  company by id
func (s *Server) UpdateCompany(c *gin.Context) {
	var cmp company
	err := c.ShouldBind(&cmp)
	if err != nil {
		s.log.Error("can not bind company", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	paramID := c.Param("id")
	id, err := uuid.Parse(paramID)
	if err != nil {
		s.log.Error("can not parse id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCompany := companyToDomain(cmp)
	err = s.companies.Update(c.Request.Context(), id, newCompany)
	if err != nil {
		s.log.Error("can not update company", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
