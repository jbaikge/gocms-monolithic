package repository

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/jbaikge/gocms/models/class"
	"github.com/jbaikge/gocms/models/document"
	"github.com/jbaikge/gocms/models/user"
	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func repositories(t *testing.T) (repos []Repository) {
	repos = make([]Repository, 0, 2)

	// Mongo Repository
	dbHost := "localhost:27017"
	if dbHostEnv := os.Getenv("DB_HOST"); dbHostEnv != "" {
		dbHost = dbHostEnv
	}

	dbName := "testing"
	if dbNameEnv := os.Getenv("DB_NAME"); dbNameEnv != "" {
		dbName = dbNameEnv
	}

	opts := options.Client().ApplyURI("mongodb://" + dbHost)

	client, err := mongo.Connect(context.Background(), opts)
	assert.NoError(t, err)

	mongoDB := client.Database(dbName)
	// Clean out db before messing with it
	assert.NoError(t, mongoDB.Drop(context.Background()))

	repos = append(repos, NewMongo(context.Background(), mongoDB))

	// Memory Repository
	repos = append(repos, NewMemory())

	return
}

func TestRepository(t *testing.T) {
	for _, repo := range repositories(t) {
		t.Run(reflect.TypeOf(repo).Elem().Name(), func(t *testing.T) {
			t.Run("DeleteClass", func(t *testing.T) {
				class := class.Class{}
				assert.NoError(t, repo.InsertClass(&class))
				assert.NoError(t, repo.DeleteClass(class.Id))
				_, err := repo.GetClassById(class.Id)
				// Make sure a "no documents in result" error pops out
				assert.Error(t, err)
			})

			t.Run("GetAllClasses", func(t *testing.T) {
				assert.NoError(t, repo.empty())
				classes := []class.Class{
					{Name: "Event", Slug: "event"},
					{Name: "Blog", Slug: "blog"},
					{Name: "News", Slug: "news"},
				}
				for _, class := range classes {
					assert.NoError(t, repo.InsertClass(&class))
				}

				all, err := repo.GetAllClasses()
				assert.NoError(t, err)
				assert.Equal(t, len(classes), len(all))

				expectSlugs := []string{"blog", "event", "news"}
				for i := range all {
					assert.Equal(t, all[i].Slug, expectSlugs[i])
				}
			})

			t.Run("GetClassById", func(t *testing.T) {
				class := class.Class{
					Slug: "get_class_by_id",
				}
				assert.NoError(t, repo.InsertClass(&class))

				check, err := repo.GetClassById(class.Id)
				assert.NoError(t, err)
				assert.Equal(t, class.Slug, check.Slug)

				_, err = repo.GetClassById(primitive.NewObjectID())
				assert.Error(t, err)
			})

			t.Run("GetClassBySlug", func(t *testing.T) {
				class := class.Class{
					Slug: "get_class_by_slug",
				}
				assert.NoError(t, repo.InsertClass(&class))

				check, err := repo.GetClassBySlug(class.Slug)
				assert.NoError(t, err)
				assert.Equal(t, class.Slug, check.Slug)

				_, err = repo.GetClassBySlug("invalid_slug")
				assert.Error(t, err)
			})

			t.Run("InsertClass", func(t *testing.T) {
				class := class.Class{
					Slug: "insert_class",
				}
				assert.NoError(t, repo.InsertClass(&class))
				assert.False(t, class.Id.IsZero())
				assert.False(t, class.Created.IsZero())
				assert.False(t, class.Updated.IsZero())
			})

			t.Run("UpdateClass", func(t *testing.T) {
				class := class.Class{
					Slug: "update_class",
				}
				assert.NoError(t, repo.InsertClass(&class))
				updated := class.Updated

				class.Slug += "_update"
				assert.NoError(t, repo.UpdateClass(&class))
				assert.True(t, updated.Before(class.Updated))

				check, err := repo.GetClassById(class.Id)
				assert.NoError(t, err)
				assert.Equal(t, class.Slug, check.Slug)

				class.Id = primitive.NewObjectID()
				class.Slug = "update_class_fail"
				assert.Error(t, repo.UpdateClass(&class))
			})

			t.Run("DeleteDocument", func(t *testing.T) {
				doc := document.Document{}
				assert.NoError(t, repo.InsertDocument(&doc))
				assert.NoError(t, repo.DeleteDocument(doc.Id))

				_, err := repo.GetDocumentById(doc.Id)
				assert.Error(t, err)
			})

			t.Run("GetChildDocumentBySlug", func(t *testing.T) {
				doc := document.Document{
					ClassId:  primitive.NewObjectID(),
					ParentId: primitive.NewObjectID(),
					Slug:     "get_child_document_slug",
				}
				assert.NoError(t, repo.InsertDocument(&doc))

				check, err := repo.GetChildDocumentBySlug(doc.ParentId, doc.Slug)
				assert.NoError(t, err)
				assert.False(t, doc.Id.IsZero())
				assert.Equal(t, doc.Id, check.Id)

				_, err = repo.GetChildDocumentBySlug(doc.ParentId, "invalid_slug")
				assert.Error(t, err)
			})

			t.Run("GetClassDocumentBySlug", func(t *testing.T) {
				doc := document.Document{
					ClassId:  primitive.NewObjectID(),
					ParentId: primitive.NewObjectID(),
					Slug:     "get_class_document_slug",
				}
				assert.NoError(t, repo.InsertDocument(&doc))

				check, err := repo.GetClassDocumentBySlug(doc.ClassId, doc.Slug)
				assert.NoError(t, err)
				assert.False(t, doc.Id.IsZero())
				assert.Equal(t, doc.Id, check.Id)

				_, err = repo.GetClassDocumentBySlug(doc.ClassId, "invalid_slug")
				assert.Error(t, err)
			})

			t.Run("GetDocumentList", func(t *testing.T) {
				classId := primitive.NewObjectID()
				ids := make([]primitive.ObjectID, 3)

				for i := range ids {
					doc := document.Document{
						ClassId: classId,
						Slug:    fmt.Sprintf("test_%d", i),
					}
					assert.NoError(t, repo.InsertDocument(&doc))
					ids[i] = doc.Id
				}

				params := document.DocumentListParams{
					ClassId: classId,
					Size:    2,
					Page:    1,
				}
				page1, err := repo.GetDocumentList(params)
				assert.NoError(t, err)
				assert.Equal(t, 3, page1.Total)
				assert.Equal(t, 2, len(page1.Documents))
				for i := range ids[0:2] {
					assert.Equal(t, ids[i], page1.Documents[i].Id)
				}

				params.Page = 2
				page2, err := repo.GetDocumentList(params)
				assert.NoError(t, err)
				assert.Equal(t, 3, page2.Total)
				assert.Equal(t, 1, len(page2.Documents))
				for i := range ids[2:3] {
					assert.Equal(t, ids[i], page1.Documents[i].Id)
				}

				params.ClassId = primitive.NewObjectID()
				noResults, err := repo.GetDocumentList(params)
				assert.NoError(t, err)
				assert.Equal(t, 0, noResults.Total)
			})

			t.Run("GetDocumentById", func(t *testing.T) {
				doc := document.Document{}
				assert.NoError(t, repo.InsertDocument(&doc))

				check, err := repo.GetDocumentById(doc.Id)
				assert.NoError(t, err)
				assert.False(t, doc.Id.IsZero())
				assert.Equal(t, doc.Id, check.Id)
			})

			t.Run("InsertDocument", func(t *testing.T) {
				doc := document.Document{
					Slug: "create_document",
				}
				assert.NoError(t, repo.InsertDocument(&doc))
				assert.False(t, doc.Id.IsZero())
				assert.False(t, doc.Created.IsZero())
				assert.False(t, doc.Updated.IsZero())
			})

			t.Run("UpdateDocument", func(t *testing.T) {
				doc := document.Document{
					ClassId: primitive.NewObjectID(),
					Slug:    "update_document",
				}
				assert.NoError(t, repo.InsertDocument(&doc))

				updated := doc.Updated
				doc.Slug += "_update"
				assert.NoError(t, repo.UpdateDocument(&doc))
				assert.True(t, updated.Before(doc.Updated))

				check, err := repo.GetDocumentById(doc.Id)
				assert.NoError(t, err)
				assert.Equal(t, doc.Slug, check.Slug)
			})

			t.Run("GetUserByEmail", func(t *testing.T) {
				u := user.User{
					Email: "test@test.com",
				}
				assert.NoError(t, repo.InsertUser(&u))

				check, err := repo.GetUserByEmail(u.Email)
				assert.NoError(t, err)
				assert.Equal(t, u.Id, check.Id)

				_, err = repo.GetUserByEmail("bad@test.com")
				assert.Error(t, err)
			})

			t.Run("GetUserById", func(t *testing.T) {
				u := user.User{}
				assert.NoError(t, repo.InsertUser(&u))

				check, err := repo.GetUserById(u.Id)
				assert.NoError(t, err)
				assert.Equal(t, u.Id, check.Id)

				_, err = repo.GetUserById(primitive.NewObjectID())
				assert.Error(t, err)
			})

			t.Run("InsertUser", func(t *testing.T) {
				u := user.User{}
				assert.NoError(t, repo.InsertUser(&u))
				assert.False(t, u.Id.IsZero())
			})

			t.Run("UpdateUser", func(t *testing.T) {
				u := user.User{
					Email: "test@test.com",
				}
				assert.NoError(t, repo.InsertUser(&u))

				u.Email = "new-email@test.com"
				assert.NoError(t, repo.UpdateUser(&u))

				check, err := repo.GetUserById(u.Id)
				assert.NoError(t, err)
				assert.Equal(t, u.Id, check.Id)

				u.Id = primitive.NewObjectID()
				assert.Error(t, repo.UpdateUser(&u))
			})
		})
	}
}
