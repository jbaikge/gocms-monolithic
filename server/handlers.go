package server

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms"
)

func (s *Server) HandleClassBuilder() gin.HandlerFunc {
	name := "admin-class-builder"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.gohtml"),
		filepath.Join(s.templatePath, "admin", "class-builder.gohtml"),
	)

	return func(c *gin.Context) {
		var class gocms.Class
		var err error

		// Pull the class by the slug to edit it
		if slug := c.Param("slug"); slug != "" {
			class, err = s.classService.GetBySlug(slug)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		if c.Request.Method == http.MethodPost {
			// Bind form values where defined
			if err := c.Bind(&class); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			// Insert or update depending on the state of class.Id
			var newUrl string
			if class.Id.IsZero() {
				newUrl = fmt.Sprintf("/admin/classes/%s/fields", class.Slug)
				err = s.classService.Insert(&class)
			} else {
				newUrl = fmt.Sprintf("/admin/classes/%s", class.Slug)
				err = s.classService.Update(&class)
			}

			// If all went well, bounce to the next page
			if err == nil {
				c.Redirect(http.StatusSeeOther, newUrl)
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
	name := "admin-class-field-builder"
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

func (s *Server) HandleClassIndex() gin.HandlerFunc {
	name := "admin-class-index"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.gohtml"),
		filepath.Join(s.templatePath, "admin", "class.gohtml"),
	)
	return func(c *gin.Context) {
		var class gocms.Class

		slug := c.Param("slug")
		if slug == "" {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		class, err := s.classService.GetBySlug(slug)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		obj := gin.H{
			"Class": class,
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
