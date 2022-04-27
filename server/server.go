package server

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms/models/class"
	"github.com/jbaikge/gocms/models/document"
)

type Server struct {
	classService    class.ClassService
	documentService document.DocumentService
	renderer        multitemplate.Renderer
	router          *gin.Engine
}

func New(
	router *gin.Engine,
	classService class.ClassService,
	documentService document.DocumentService,
) *Server {
	renderer := multitemplate.NewRenderer()
	router.HTMLRender = renderer
	return &Server{
		classService:    classService,
		documentService: documentService,
		renderer:        renderer,
		router:          router,
	}
}

func (s Server) Run(listenAddress string) error {
	routes := s.Routes()
	return routes.Run(listenAddress)
}
