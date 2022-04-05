package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms"
)

func (s *Server) HandleClassBuilder() gin.HandlerFunc {
	name := "class-builder"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.gohtml"),
		filepath.Join(s.templatePath, "admin", "class-builder.gohtml"),
	)

	return func(c *gin.Context) {
		var class gocms.Class
		var err error

		if err := c.Bind(&class); err != nil {
			log.Print(err)
		}
		if c.Request.Method == http.MethodPost {
			err = s.classService.Insert(&class)
			if err == nil {
				newUrl := fmt.Sprintf("/admin/classes/%s/fields", class.Slug)
				c.Redirect(http.StatusTemporaryRedirect, newUrl)
				return
			}
		}

		obj := gin.H{
			"Class": class,
			"Error": err,
		}
		if list, ok := c.Get("classList"); ok {
			obj["ClassList"] = list
		}
		c.HTML(http.StatusOK, name, obj)
	}
}

func (s *Server) HandleClassFieldBuilder() gin.HandlerFunc {
	name := "class-field-builder"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.gohtml"),
		filepath.Join(s.templatePath, "admin", "class-field-builder.gohtml"),
	)

	types := []struct {
		Type     string `json:"type"`
		Label    string `json:"label"`
		Template string `json:"template"`
	}{
		{gocms.TypeDate, "Date", "date"},
		{gocms.TypeDateTime, "Date & Time", "date"},
		{gocms.TypeEmail, "Email", "email"},
		{gocms.TypeMultiSelect, "Multi-Select", "select"},
		{gocms.TypeNumber, "Number", "number"},
		{gocms.TypeSelect, "Select", "select"},
		{gocms.TypeText, "Text", "text"},
		{gocms.TypeTextArea, "Textarea", "textarea"},
		{gocms.TypeTime, "Time", "time"},
		{gocms.TypeTinyMCE, "TinyMCE", "tinymce"},
		{gocms.TypeUpload, "Upload", "upload"},
	}

	return func(c *gin.Context) {
		var class gocms.Class
		obj := gin.H{
			"FieldTypes": types,
			"Class":      class,
		}
		if list, ok := c.Get("classList"); ok {
			obj["ClassList"] = list
		}
		c.HTML(http.StatusOK, name, obj)
	}
}

func (s *Server) HandleNavBar() gin.HandlerFunc {
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
