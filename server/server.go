package server

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms"
)

type Server struct {
	classService    gocms.ClassService
	documentService gocms.DocumentService
	renderer        multitemplate.Renderer
	router          *gin.Engine
	templatePath    string
}

func New(
	templatePath string,
	router *gin.Engine,
	classService gocms.ClassService,
	documentService gocms.DocumentService,
) *Server {
	renderer := multitemplate.NewRenderer()
	router.HTMLRender = renderer
	return &Server{
		classService:    classService,
		documentService: documentService,
		renderer:        renderer,
		router:          router,
		templatePath:    templatePath,
	}
}

func (s Server) Run(listenAddress string) error {
	routes := s.Routes()
	return routes.Run(listenAddress)
}
