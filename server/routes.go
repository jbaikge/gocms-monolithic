package server

import (
	"crypto/rand"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.router

	authKey := make([]byte, 64)
	rand.Read(authKey)
	store := cookie.NewStore(authKey)
	router.Use(sessions.Sessions("gocms", store))

	// router.GET("/assets/:filename", s.HandleAsset())
	// router.GET("/forms/:id", s.HandleForm())
	// router.POST("/forms/:id", s.HandleForm())

	admin := router.Group("/admin")
	admin.Use(s.MiddlewareAdminAuth())
	admin.Use(s.MiddlewareNavBar())
	{
		admin.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, "/admin/login")
		})

		admin.GET("/login", s.HandleAdminLogin())
		admin.POST("/login", s.HandleAdminLogin())

		// Send successful logins here. Maybe just display a table with classes
		// and document counts?
		// admin.GET("/dashboard", s.HandleAdminDashboard())

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
