package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms"
	"github.com/jbaikge/gocms/repository"
	"github.com/zeebo/assert"
)

func TestGetContext(t *testing.T) {
	var i int
	var f float64
	c := new(gin.Context)

	c.Set("myInt", 42)

	assert.NoError(t, getContext(c, "myInt", &i))
	assert.Error(t, getContext(c, "myInt", &f))
	assert.Error(t, getContext(c, "myFloat", &f))
}

func TestServer(t *testing.T) {
	router := gin.Default()
	repo := repository.NewMemory()
	classRepository := gocms.NewClassService(repo)
	docRepository := gocms.NewDocumentService(repo)
	s := New(router, classRepository, docRepository)
	routes := s.Routes()

	t.Run("MiddlewareClass", func(t *testing.T) {
		// Requesting a class that doesn't exist should bounce back with a 404.
		t.Run("BadClass", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin/classes/bad/edit", nil)
			w := httptest.NewRecorder()
			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		// An empty class slug should bounce back with a 404
		t.Run("NoClass", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin/classes//edit", nil)
			w := httptest.NewRecorder()
			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		// A good class will fall through to the correct handler
		t.Run("GoodClass", func(t *testing.T) {
			class := gocms.Class{Slug: "good_class"}
			assert.NoError(t, repo.InsertClass(&class))
			req := httptest.NewRequest(http.MethodGet, "/admin/classes/good_class/edit", nil)
			w := httptest.NewRecorder()
			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("HandleClassBuilder", func(t *testing.T) {
		// Not a whole lot happens here - the new class form displays
		t.Run("GetNew", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin/classes/new", nil)
			w := httptest.NewRecorder()
			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		// Empty POST should basically refresh the page with Errors available
		// to show what went wrong.
		t.Run("PostNewEmpty", func(t *testing.T) {
			var values url.Values
			body := strings.NewReader(values.Encode())
			req := httptest.NewRequest(http.MethodPost, "/admin/classes/new", body)
			w := httptest.NewRecorder()
			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})

		// Valid POST will bounce to the field builder page with a 303 - See
		// Other response
		t.Run("PostNewClass", func(t *testing.T) {
			values := make(url.Values)
			values.Set("name", "Blogs")
			values.Set("slug", "blog")
			values.Set("singular_name", "Blog")
			values.Set("menu_label", "Blogs")
			values.Set("add_item_label", "Add Blog Entry")
			values.Set("new_item_label", "New Blog Entry")
			values.Set("edit_item_label", "Edit Blog Entry")
			values.Set("table_labels", "Title Slug")
			values.Set("table_fields", "title slug")
			body := strings.NewReader(values.Encode())

			req := httptest.NewRequest(http.MethodPost, "/admin/classes/new", body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusSeeOther, w.Code)
		})
	})
}
