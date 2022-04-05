package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// router.GET("/assets/:filename", s.HandleAsset())
	// router.GET("/forms/:id", s.HandleForm())
	// router.POST("/forms/:id", s.HandleForm())

	admin := router.Group("/admin")
	admin.Use(s.HandleNavBar())
	{
		classes := admin.Group("/classes")
		{
			classes.GET("/new", s.HandleClassBuilder())
			classes.POST("/new", s.HandleClassBuilder())
			class := classes.Group("/:slug")
			{
				// class.GET("/", s.HandleClassIndex())
				class.GET("/edit", s.HandleClassBuilder())
				class.POST("/edit", s.HandleClassBuilder())
				class.GET("/fields", s.HandleClassFieldBuilder())
				class.POST("/fields", s.HandleClassFieldBuilder())
			}
		}
		// forms := admin.Group("/forms")
		// 	forms.GET("/new", s.HandleFormBuilder())
		// 	forms.POST("/new", s.HandleFormBuilder())
		// }
	}

	return router
}
