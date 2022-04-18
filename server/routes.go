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
	admin.Use(s.MiddlewareNavBar())
	{
		classes := admin.Group("/classes")
		{
			classes.GET("/new", s.HandleClassBuilder())
			classes.POST("/new", s.HandleClassBuilder())
			class := classes.Group("/:class")
			class.Use(s.MiddlewareClass())
			{
				class.GET("/", s.HandleDocumentList())
				class.GET("/edit", s.HandleClassBuilder())
				class.POST("/edit", s.HandleClassBuilder())
				class.GET("/fields", s.HandleClassFieldBuilderGet())
				class.POST("/fields", s.HandleClassFieldBuilderPost())
				class.GET("/new", s.HandleDocumentBuilder())
				class.POST("/new", s.HandleDocumentBuilder())
				class.GET("/:doc_id", s.HandleDocumentBuilder())
				class.POST("/:doc_id", s.HandleDocumentBuilder())
			}
		}
		// forms := admin.Group("/forms")
		// 	forms.GET("/new", s.HandleFormBuilder())
		// 	forms.POST("/new", s.HandleFormBuilder())
		// }
	}

	return router
}
