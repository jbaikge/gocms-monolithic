package server

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms/models/class"
	"github.com/jbaikge/gocms/models/document"
	"github.com/jbaikge/gocms/models/user"
)

type Server struct {
	classService    class.ClassService
	documentService document.DocumentService
	userService     user.UserService
	renderer        multitemplate.Renderer
	router          *gin.Engine
}

func New(
	router *gin.Engine,
	classService class.ClassService,
	documentService document.DocumentService,
	userService user.UserService,
) *Server {
	renderer := multitemplate.NewRenderer()
	router.HTMLRender = renderer
	return &Server{
		classService:    classService,
		documentService: documentService,
		userService:     userService,
		renderer:        renderer,
		router:          router,
	}
}

func (s Server) Run(listenAddress string) error {
	routes := s.Routes()
	return routes.Run(listenAddress)
}
