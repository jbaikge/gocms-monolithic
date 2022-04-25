package server

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:embed templates
var fs embed.FS

func getContext[T any](c *gin.Context, key string, into *T) (err error) {
	if obj, ok := c.Get(key); ok {
		if t, ok := obj.(T); ok {
			*into = t
			return
		}
		return fmt.Errorf("could not cast value")
	}
	return fmt.Errorf("key not found: %s", key)
}

func (s *Server) HandleClassBuilder() gin.HandlerFunc {
	name := "admin-class-builder"
	s.renderer.Add(name, template.Must(template.New("base.html").ParseFS(
		fs,
		"templates/admin/base.html",
		"templates/admin/class-builder.html",
	)))

	return func(c *gin.Context) {
		var class gocms.Class
		var err error

		// If no Class, then we are on /new
		if _, ok := c.Get("class"); ok {
			// Class will be set by the middleware preceding this handler
			_ = getContext(c, "class", &class)
		}

		if c.Request.Method == http.MethodPost {
			// Bind form values where defined
			if err := c.Bind(&class); err != nil {
				return
			}

			// Insert or update depending on the state of class.Id
			var newUrl string
			if class.Id.IsZero() {
				newUrl = fmt.Sprintf("/admin/classes/%s/fields", class.Slug)
				err = s.classService.Insert(&class)
			} else {
				newUrl = fmt.Sprintf("/admin/classes/%s/", class.Slug)
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
	s.renderer.Add(name, template.Must(template.New("base.html").ParseFS(
		fs,
		"templates/admin/base.html",
		"templates/admin/class-field-builder.html",
	)))

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
		{gocms.TypeSelect, "Select (Class)", "select-class"},
		{gocms.TypeSelect, "Select (Static)", "select-static"},
		{gocms.TypeText, "Text", "text"},
		{gocms.TypeTextArea, "Textarea", "textarea"},
		{gocms.TypeTime, "Time", "time"},
		{gocms.TypeTinyMCE, "TinyMCE", "tinymce"},
		{gocms.TypeUpload, "Upload", "upload"},
	}

	return func(c *gin.Context) {
		obj := gin.H{
			"FieldTypes": types,
			"Error":      nil,
		}
		if list, ok := c.Get("classList"); ok {
			obj["ClassList"] = list
		}
		if class, ok := c.Get("class"); ok {
			obj["Class"] = class
		}

		c.HTML(http.StatusOK, name, obj)
	}
}

func (s *Server) HandleClassFieldBuilderPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var class gocms.Class

		// class gauranteed to be set per middleware preceding this handler
		_ = getContext(c, "class", &class)

		if err := c.Bind(&class); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

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

func (s *Server) HandleDocumentBuilder() gin.HandlerFunc {
	name := "admin-document-builder"
	s.renderer.Add(name, template.Must(template.New("base.html").ParseFS(
		fs,
		"templates/admin/base.html",
		"templates/admin/document-builder.html",
	)))

	layout := "2006-01-02T15:04"
	loc, _ := time.LoadLocation("America/New_York")

	return func(c *gin.Context) {
		var class gocms.Class
		var doc gocms.Document

		// Class gauranteed to be set from middleware preceding this handler
		_ = getContext(c, "class", &class)

		if id := c.Param("doc_id"); id != "" {
			bsonId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			if doc, err = s.documentService.GetById(bsonId); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else {
			doc.ClassId = class.Id
			doc.Published = time.Now().In(loc)
		}

		if c.Request.Method == http.MethodPost {
			doc.Title = c.PostForm("title")
			doc.Slug = c.PostForm("slug")
			if published, err := time.ParseInLocation(layout, c.PostForm("published"), loc); err == nil {
				doc.Published = published
			}
			if doc.Values == nil {
				doc.Values = make(map[string]interface{})
			}
			for _, field := range class.Fields {
				doc.Values[field.Name] = c.PostForm(field.Name)
			}
			if doc.Id.IsZero() {
				if err := s.documentService.Insert(&doc); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
			} else {
				if err := s.documentService.Update(&doc); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
			}

			c.Redirect(http.StatusSeeOther, "/admin/classes/"+class.Slug+"/"+doc.Id.Hex())
			return
		}

		for i, field := range class.Fields {
			if field.DataSourceId.IsZero() {
				continue
			}
			docs := []gocms.Document{
				{
					Id:    primitive.NewObjectID(),
					Slug:  "moo",
					Title: "Cow",
				},
			}
			for _, doc := range docs {
				class.Fields[i].Options += fmt.Sprintf(
					"%s|%s\n",
					field.Apply(doc.Value(field.DataSourceValue)),
					field.Apply(doc.Value(field.DataSourceLabel)),
				)
			}
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

func (s *Server) HandleDocumentList() gin.HandlerFunc {
	name := "admin-document-list"
	s.renderer.Add(name, template.Must(template.New("base.html").ParseFS(
		fs,
		"templates/admin/base.html",
		"templates/admin/document-list.html",
	)))

	return func(c *gin.Context) {
		var class gocms.Class

		// Class gauranteed to be set by middleware preceding this handler
		_ = getContext(c, "class", &class)

		page, err := strconv.ParseInt(c.Query("p"), 10, 64)
		if err != nil || page == 0 {
			page = 1
		}

		params := gocms.DocumentListParams{
			ClassId: class.Id,
			Page:    page,
			Size:    2,
		}
		list, err := s.documentService.List(params)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		table := NewTable(class, list.Documents)

		obj := gin.H{
			"Class":      class,
			"Table":      table,
			"Pagination": NewPagination(params.Page, params.Size, list.Total),
		}
		if list, ok := c.Get("classList"); ok {
			obj["ClassList"] = list
		}

		c.HTML(http.StatusOK, name, obj)
	}
}
