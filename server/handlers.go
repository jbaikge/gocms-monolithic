package server

import (
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

	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, name, gin.H{})
	}
}

func (s *Server) HandleClassFieldBuilder() gin.HandlerFunc {
	name := "class-field-builder"
	s.renderer.AddFromFiles(
		name,
		filepath.Join(s.templatePath, "admin", "base.gohtml"),
		filepath.Join(s.templatePath, "admin", "class-field-builder.gohtml"),
	)

	type fieldType struct {
		Type     string `json:"type"`
		Label    string `json:"label"`
		Template string `json:"template"`
	}
	types := []fieldType{
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

	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, name, gin.H{
			"FieldTypes":  types,
			"ClassFields": []gocms.Field{},
		})
	}
}
