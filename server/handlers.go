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
		{
			Type:     gocms.TypeDate,
			Label:    "Date",
			Template: "date",
		},
		{
			Type:     gocms.TypeDateTime,
			Template: "date",
			Label:    "Date & Time",
		},
		{
			Type:     gocms.TypeEmail,
			Template: "email",
			Label:    "Email",
		},
		{
			Type:     gocms.TypeMultiSelect,
			Template: "select",
			Label:    "Multi-Select",
		},
		{
			Type:     gocms.TypeNumber,
			Template: "number",
			Label:    "Number",
		},
		{
			Type:     gocms.TypeSelect,
			Template: "select",
			Label:    "Select",
		},
		{
			Type:     gocms.TypeText,
			Template: "text",
			Label:    "Text",
		},
		{
			Type:     gocms.TypeTextArea,
			Template: "textarea",
			Label:    "Textarea",
		},
		{
			Type:     gocms.TypeTime,
			Template: "time",
			Label:    "Time",
		},
		{
			Type:     gocms.TypeTinyMCE,
			Template: "tinymce",
			Label:    "TinyMCE",
		},
		{
			Type:     gocms.TypeUpload,
			Template: "upload",
			Label:    "Upload",
		},
	}

	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, name, gin.H{
			"FieldTypes":  types,
			"ClassFields": []gocms.Field{},
		})
	}
}
