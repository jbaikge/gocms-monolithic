package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.router

	admin := router.Group("/admin")
	{
		classes := admin.Group("/classes")
		{
			classes.GET("/new", s.HandleClassBuilder())
			classes.POST("/new", s.HandleClassBuilder())
			class := classes.Group("/:slug")
			{
				class.GET("/edit", s.HandleClassBuilder())
				class.POST("/edit", s.HandleClassBuilder())
				class.GET("/fields", s.HandleClassFieldBuilder())
				class.POST("/fields", s.HandleClassFieldBuilder())
			}
		}
	}

	return router
}
