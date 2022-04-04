package server

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

type Server struct {
	renderer     multitemplate.Renderer
	router       *gin.Engine
	templatePath string
}

func New(templatePath string, router *gin.Engine) *Server {
	renderer := multitemplate.NewRenderer()
	router.HTMLRender = renderer
	return &Server{
		renderer:     renderer,
		router:       router,
		templatePath: templatePath,
	}
}

func (s Server) Run(listenAddress string) error {
	routes := s.Routes()
	return routes.Run(listenAddress)
}
