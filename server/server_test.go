package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms"
	"github.com/jbaikge/gocms/repository"
	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	classService := gocms.NewClassService(repo)
	docService := gocms.NewDocumentService(repo)
	s := New(router, classService, docService)
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

			uri := "/admin/classes/new"
			req := httptest.NewRequest(http.MethodPost, uri, body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusSeeOther, w.Code)

			location, err := w.Result().Location()
			assert.NoError(t, err)
			assert.Equal(t, "/admin/classes/blog/fields", location.Path)
		})

		// Invalid POST should return a BadRequest
		t.Run("PostBadClass", func(t *testing.T) {
			values := make(url.Values)
			values.Add("parents", "moo") // Non array and not a BSON ID
			body := strings.NewReader(values.Encode())

			uri := "/admin/classes/new"
			req := httptest.NewRequest(http.MethodPost, uri, body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			w := httptest.NewRecorder()

			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		// Editing an existing class should redirect to the listing page after
		// POST
		t.Run("PostClassUpdate", func(t *testing.T) {
			values := make(url.Values)
			values.Set("name", "Objects")
			values.Set("slug", "objects")
			body := strings.NewReader(values.Encode())

			// Push class in
			uri := "/admin/classes/new"
			req := httptest.NewRequest(http.MethodPost, uri, body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusSeeOther, w.Code)

			// Update class and check if we bounce to the listing page
			uri = "/admin/classes/objects/edit"
			body = strings.NewReader(values.Encode())
			req = httptest.NewRequest(http.MethodPost, uri, body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w = httptest.NewRecorder()
			routes.ServeHTTP(w, req)
			assert.Equal(t, http.StatusSeeOther, w.Code)

			location, err := w.Result().Location()
			assert.NoError(t, err)
			assert.Equal(t, "/admin/classes/objects/", location.Path)
		})
	})

	t.Run("HandleClassFieldBuilderGet", func(t *testing.T) {
		// Nothing really exciting happens, a bunch of data gets pushed into the
		// template context and rendered in the browser
		class := gocms.Class{Slug: "builder_get"}
		assert.NoError(t, repo.InsertClass(&class))
		target := "/admin/classes/" + class.Slug + "/fields"
		req := httptest.NewRequest(http.MethodGet, target, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("HandleClassFieldsBuilderPost", func(t *testing.T) {
		type request struct {
			Fields []gocms.Field `json:"fields"`
		}

		type response struct {
			Success bool   `json:"success"`
			Error   string `json:"error"`
		}

		// Follow all the rules
		t.Run("Good", func(t *testing.T) {
			class := gocms.Class{Name: "Builder Good", Slug: "builder_good"}
			assert.NoError(t, repo.InsertClass(&class))

			var data request
			data.Fields = []gocms.Field{
				{
					Name:  "my_field",
					Label: "My Field",
					Type:  gocms.TypeText,
				},
			}
			body := new(bytes.Buffer)
			assert.NoError(t, json.NewEncoder(body).Encode(data))

			target := "/admin/classes/" + class.Slug + "/fields"
			req := httptest.NewRequest(http.MethodPost, target, body)
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusAccepted, w.Code)

			var resp response
			assert.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			assert.True(t, resp.Success)
			assert.Equal(t, "", resp.Error)
		})

		// Trigger ClassService.Validate to fail
		t.Run("Bad", func(t *testing.T) {
			class := gocms.Class{Name: "Builder Bad", Slug: "builder_bad"}
			assert.NoError(t, repo.InsertClass(&class))

			var data request
			data.Fields = []gocms.Field{
				{
					Name:  "my_field",
					Label: "My Label",
				},
			}
			body := new(bytes.Buffer)
			assert.NoError(t, json.NewEncoder(body).Encode(data))

			target := "/admin/classes/" + class.Slug + "/fields"
			req := httptest.NewRequest(http.MethodPost, target, body)
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var resp response
			assert.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			assert.False(t, resp.Success)
			assert.Equal(t, "field[0] type is empty", resp.Error)
		})

		// Sending in field data as a non-array
		t.Run("Malformed", func(t *testing.T) {
			class := gocms.Class{Name: "Builder Malformed", Slug: "builder_malformed"}
			assert.NoError(t, repo.InsertClass(&class))

			data := struct {
				Fields gocms.Field
			}{
				Fields: gocms.Field{
					Name:  "my_field",
					Label: "My Field",
					Type:  gocms.TypeText,
				},
			}
			body := new(bytes.Buffer)
			assert.NoError(t, json.NewEncoder(body).Encode(data))

			target := "/admin/classes/" + class.Slug + "/fields"
			req := httptest.NewRequest(http.MethodPost, target, body)
			req.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var resp response
			assert.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			assert.False(t, resp.Success)
			assert.Equal(t, "json: cannot unmarshal", resp.Error[:22])
		})
	})

	t.Run("HandleDocumentBuilder", func(t *testing.T) {
		type response struct {
			Document gocms.Document
			Class    gocms.Class
			Error    error
		}
		class := gocms.Class{
			Name: "Doc Builder",
			Slug: "builder_class",
			Fields: []gocms.Field{
				{
					Name:  "field_1",
					Label: "Field 1",
					Type:  gocms.TypeText,
				},
			},
		}
		assert.NoError(t, repo.InsertClass(&class))
		baseURL := "/admin/classes/" + class.Slug

		doc := gocms.Document{ClassId: class.Id, Slug: "builder_doc", Title: "Doc"}
		assert.NoError(t, repo.InsertDocument(&doc))

		// Retrieve the details for a new document form, ask for JSON instead
		// of trying to parse the HTML
		t.Run("New Form", func(t *testing.T) {
			target := baseURL + "/new"
			req := httptest.NewRequest(http.MethodGet, target, nil)
			req.Header.Add("Accept", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp response
			assert.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
			assert.False(t, resp.Document.ClassId.IsZero())
			assert.False(t, resp.Document.Published.IsZero())
		})

		t.Run("Post New", func(t *testing.T) {
			values := make(url.Values)
			values.Set("title", "Post Document")
			values.Set("slug", "post_document")
			values.Set("published", time.Now().Format("2006-01-02T15:04"))
			values.Set("field_1", "value_1")
			values.Set("field_2", "value_2") // Extra, should be ignored
			body := strings.NewReader(values.Encode())

			target := baseURL + "/new"
			req := httptest.NewRequest(http.MethodPost, target, body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// For now assert there's a bounce, might figure out how to verify
			// the URL later
			assert.Equal(t, http.StatusSeeOther, w.Code)
		})

		t.Run("Update", func(t *testing.T) {
			values := make(url.Values)
			values.Set("title", "Updated Document")
			values.Set("slug", doc.Slug)
			values.Set("field_1", "Updated Field 1")
			body := strings.NewReader(values.Encode())

			target := baseURL + "/" + doc.Id.Hex()
			req := httptest.NewRequest(http.MethodPost, target, body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusSeeOther, w.Code)

			// Essentially refresh the page
			location, _ := w.Result().Location()
			assert.Equal(t, location.Path, target)
		})

		t.Run("Fail Validation", func(t *testing.T) {
			values := make(url.Values)
			body := strings.NewReader(values.Encode())

			target := baseURL + "/new"
			req := httptest.NewRequest(http.MethodPost, target, body)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			// Insert with missing data
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)

			// Update with missing data
			req.URL.Path = baseURL + "/" + doc.Id.Hex()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("Bad ID", func(t *testing.T) {
			target := baseURL + "/lol"
			req := httptest.NewRequest(http.MethodGet, target, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Code)

			req.URL.Path = baseURL + "/" + primitive.NewObjectID().Hex()
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	})

	t.Run("HandleDocumentList", func(t *testing.T) {
		class := gocms.Class{Name: "Doc Test", Slug: "doc_test"}
		assert.NoError(t, repo.InsertClass(&class))
		baseURL := "/admin/classes/" + class.Slug + "/"

		numDocs := 10
		for i := 0; i < numDocs; i++ {
			doc := gocms.Document{
				ClassId: class.Id,
				Title:   fmt.Sprintf("Document %d", i),
				Slug:    fmt.Sprintf("doc_%d", i),
			}
			assert.NoError(t, repo.InsertDocument(&doc))
		}

		// Default listing landing page
		t.Run("Landing", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, baseURL, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			body := w.Body.String()
			assert.Equal(t, numDocs+1, strings.Count(body, "<tr>"))
		})

		// Change the number of items listed per page
		t.Run("Per Page", func(t *testing.T) {
			target := baseURL + "?pp=5"
			req := httptest.NewRequest(http.MethodGet, target, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			body := w.Body.String()
			assert.Equal(t, 5+1, strings.Count(body, "<tr>"))
		})

		// Confirm pagination works
		t.Run("Pagination", func(t *testing.T) {
			target := baseURL + "?pp=1&p=3"
			req := httptest.NewRequest(http.MethodGet, target, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			body := w.Body.String()
			assert.True(t, strings.Contains(body, "Document 2"))
		})
	})
}
