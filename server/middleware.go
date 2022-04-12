package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) MiddlewareClassInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("class")
		if slug == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		class, err := s.classService.GetBySlug(slug)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.Set("class", class)
		c.Next()
	}
}

func (s *Server) MiddlewareNavBar() gin.HandlerFunc {
	return func(c *gin.Context) {
		all, err := s.classService.All()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Set("classList", all)
		c.Next()
	}
}
