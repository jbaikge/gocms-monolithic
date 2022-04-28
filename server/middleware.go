package server

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) MiddlewareAdminAuth() gin.HandlerFunc {
	path := "/admin/login"

	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("adminUserId")
		switch value := v.(type) {
		case primitive.ObjectID:
			log.Printf("User ID: %s", value.Hex())
			c.Next()
		default:
			if c.Request.URL.Path == path {
				// Prevent an infinite loop
				return
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, path)
		c.Abort()
	}
}

func (s *Server) MiddlewareClass() gin.HandlerFunc {
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
