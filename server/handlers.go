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
		filepath.Join(s.templatePath, "admin", "base.html"),
		filepath.Join(s.templatePath, "admin", "class-builder.html"),
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

func (s *Server) HandleClassFieldBuilderGet() gin.HandlerFunc {
	name := "admin-class-field-builder"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.html"),
		filepath.Join(s.templatePath, "admin", "class-field-builder.html"),
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

	type postData struct {
		Fields []gocms.Field `form:"fields" json:"fields"`
	}

	return func(c *gin.Context) {
		var class gocms.Class
		var err error

		if class, err = s.classService.GetBySlug(c.Param("slug")); err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		obj := gin.H{
			"FieldTypes": types,
			"Class":      class,
			"Error":      err,
		}
		if list, ok := c.Get("classList"); ok {
			obj["ClassList"] = list
		}
		c.HTML(http.StatusOK, name, obj)
	}
}

func (s *Server) HandleClassFieldBuilderPost() gin.HandlerFunc {
	type postData struct {
		Fields []gocms.Field
	}

	return func(c *gin.Context) {
		var class gocms.Class
		slug := c.Param("slug")
		class, err := s.classService.GetBySlug(slug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Class with slug, " + slug + ", not found",
			})
			return
		}

		var post postData
		if err := c.Bind(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		class.Fields = post.Fields
		if err := s.classService.Update(&class); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"success": true})
	}
}

func (s *Server) HandleClassIndex() gin.HandlerFunc {
	name := "admin-class-index"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.html"),
		filepath.Join(s.templatePath, "admin", "class.html"),
	)
	return func(c *gin.Context) {
		var class gocms.Class

		slug := c.Param("slug")
		if slug == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		class, err := s.classService.GetBySlug(slug)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
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

func (s *Server) HandleDocumentBuilder() gin.HandlerFunc {
	name := "admin-document-builder"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.html"),
		filepath.Join(s.templatePath, "admin", "document-builder.html"),
	)
	return func(c *gin.Context) {
		var class gocms.Class
		var doc gocms.Document

		slug := c.Param("slug")
		if slug == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		class, err := s.classService.GetBySlug(slug)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		obj := gin.H{
			"Document": doc,
			"Class":    class,
			"Error":    nil,
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
