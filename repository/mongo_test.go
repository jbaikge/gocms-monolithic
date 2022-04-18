package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jbaikge/gocms"
	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func XXXTestMongo(t *testing.T) {
	const classColl = "classes"
	const docColl = "documents"

	dbHost := "localhost:27017"
	if dbHostEnv := os.Getenv("DB_HOST"); dbHostEnv != "" {
		dbHost = dbHostEnv
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+dbHost))
	assert.NoError(t, err)

	db := client.Database("testing")
	// Clean out db before messing with it
	assert.NoError(t, db.Drop(ctx))

	repo := NewMongo(ctx, db)

	t.Run("Class", func(t *testing.T) {
		t.Run("All", func(t *testing.T) {
			// need a clean collection
			assert.NoError(t, db.Collection("classes").Drop(ctx))
			classes := []gocms.Class{
				{Name: "Event", Slug: "event"},
				{Name: "Blog", Slug: "blog"},
				{Name: "News", Slug: "news"},
			}
			for _, class := range classes {
				assert.NoError(t, repo.InsertClass(&class))
			}

			all, err := repo.GetAllClasses()
			assert.NoError(t, err)
			assert.Equal(t, len(all), len(classes))

			expectSlugs := []string{"blog", "event", "news"}
			for i := range all {
				assert.Equal(t, all[i].Slug, expectSlugs[i])
			}
		})

		t.Run("Create", func(t *testing.T) {
			class := gocms.Class{
				Slug: "create_test",
			}
			assert.NoError(t, repo.InsertClass(&class))
			assert.False(t, class.Id.IsZero())
			assert.False(t, class.Created.IsZero())
			assert.False(t, class.Updated.IsZero())
		})

		t.Run("Read", func(t *testing.T) {
			class := gocms.Class{
				Slug: "read_test",
			}
			assert.NoError(t, repo.InsertClass(&class))

			idCheck, err := repo.GetClassById(class.Id)
			assert.NoError(t, err)
			assert.Equal(t, class.Slug, idCheck.Slug)

			slugCheck, err := repo.GetClassBySlug(class.Slug)
			assert.NoError(t, err)
			assert.Equal(t, class.Slug, slugCheck.Slug)
		})

		t.Run("Update", func(t *testing.T) {
			class := gocms.Class{
				Slug: "update_test",
			}
			assert.NoError(t, repo.InsertClass(&class))
			updated := class.Updated

			class.Slug = "update_test_update"
			assert.NoError(t, repo.UpdateClass(&class))
			assert.True(t, updated.Before(class.Updated))

			check, err := repo.GetClassById(class.Id)
			assert.NoError(t, err)
			assert.Equal(t, class.Slug, check.Slug)
		})

		t.Run("Update Not Found", func(t *testing.T) {
			class := gocms.Class{
				Slug: "update_not_found",
			}
			assert.NoError(t, repo.InsertClass(&class))

			class.Id = primitive.NewObjectID()
			class.Slug = "should_not_work"
			assert.Error(t, repo.UpdateClass(&class))
		})

		t.Run("Delete", func(t *testing.T) {
			class := gocms.Class{}
			assert.NoError(t, repo.InsertClass(&class))
			assert.NoError(t, repo.DeleteClass(class.Id))
			_, err := repo.GetClassById(class.Id)
			// Make sure a "no documents in result" error pops out
			assert.Error(t, err)
		})
	})

	t.Run("Document", func(t *testing.T) {
		// Class used for tests below
		class := gocms.Class{
			Slug: "class_test",
		}
		assert.NoError(t, repo.InsertClass(&class))

		t.Run("Create", func(t *testing.T) {
			doc := gocms.Document{
				Slug: "create_test",
			}
			assert.NoError(t, repo.InsertDocument(&doc))
			assert.False(t, doc.Id.IsZero())
			assert.False(t, doc.Created.IsZero())
			assert.False(t, doc.Updated.IsZero())
		})

		t.Run("Read", func(t *testing.T) {
			doc := gocms.Document{
				ClassId:  class.Id,
				ParentId: primitive.NewObjectID(),
				Slug:     "read_test",
			}
			assert.NoError(t, repo.InsertDocument(&doc))

			t.Run("GetDocumentById", func(t *testing.T) {
				check, err := repo.GetDocumentById(doc.Id)
				assert.NoError(t, err)
				assert.Equal(t, doc.ClassId, check.ClassId)
			})

			t.Run("GetChildDocumentBySlug", func(t *testing.T) {
				check, err := repo.GetChildDocumentBySlug(doc.ParentId, doc.Slug)
				assert.NoError(t, err)
				assert.Equal(t, doc.Id, check.Id)
			})

			t.Run("GetClassDocumentBySlug", func(t *testing.T) {
				check, err := repo.GetClassDocumentBySlug(doc.ClassId, doc.Slug)
				assert.NoError(t, err)
				assert.Equal(t, doc.Id, check.Id)
			})
		})

		t.Run("Update", func(t *testing.T) {
			doc := gocms.Document{
				ClassId: class.Id,
				Slug:    "update_test",
			}
			assert.NoError(t, repo.InsertDocument(&doc))

			updated := doc.Updated
			doc.Slug = "update_test_update"
			assert.NoError(t, repo.UpdateDocument(&doc))
			assert.True(t, updated.Before(doc.Updated))

			check, err := repo.GetDocumentById(doc.Id)
			assert.NoError(t, err)
			assert.Equal(t, doc.Slug, check.Slug)
		})

		t.Run("Delete", func(t *testing.T) {
			doc := gocms.Document{}
			assert.NoError(t, repo.InsertDocument(&doc))
			assert.NoError(t, repo.DeleteDocument(doc.Id))

			_, err := repo.GetDocumentById(doc.Id)
			assert.Error(t, err)
		})

		t.Run("List", func(t *testing.T) {
			classId := primitive.NewObjectID()
			ids := make([]primitive.ObjectID, 3)

			for i := range ids {
				doc := gocms.Document{
					ClassId: classId,
					Slug:    fmt.Sprintf("test_%d", i),
				}
				assert.NoError(t, repo.InsertDocument(&doc))
				ids[i] = doc.Id
			}

			params := gocms.DocumentListParams{
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
	})
}
