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
		{gocms.TypeDateTime, "date", "Date & Time"},
		{gocms.TypeEmail, "email", "Email"},
		{gocms.TypeMultiSelect, "select", "Multi-Select"},
		{gocms.TypeNumber, "number", "Number"},
		{gocms.TypeSelect, "select", "Select"},
		{gocms.TypeText, "text", "Text"},
		{gocms.TypeTextArea, "textarea", "Textarea"},
		{gocms.TypeTime, "time", "Time"},
		{gocms.TypeTinyMCE, "tinymce", "TinyMCE"},
		{gocms.TypeUpload, "upload", "Upload"},
	}

	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, name, gin.H{
			"FieldTypes":  types,
			"ClassFields": []gocms.Field{},
		})
	}
}
