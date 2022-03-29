package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jbaikge/gocms/pkg/model"
	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongo(t *testing.T) {
	const classColl = "classes"
	const docColl = "documents"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err)

	db := client.Database("testing")
	t.Cleanup(func() {
		db.Drop(ctx)
	})

	repo := NewMongo(ctx, db)

	t.Run("Class", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {
			class := model.Class{
				Slug: "create_test",
			}
			assert.NoError(t, repo.InsertClass(&class))
			assert.False(t, class.Id.IsZero())
			assert.False(t, class.Created.IsZero())
			assert.False(t, class.Updated.IsZero())
		})

		t.Run("Read", func(t *testing.T) {
			class := model.Class{
				Slug: "read_test",
			}
			assert.NoError(t, repo.InsertClass(&class))

			check, err := repo.GetClass(class.Id)
			assert.NoError(t, err)
			assert.Equal(t, class.Slug, check.Slug)
		})

		t.Run("Update", func(t *testing.T) {
			class := model.Class{
				Slug: "update_test",
			}
			assert.NoError(t, repo.InsertClass(&class))
			updated := class.Updated

			class.Slug = "update_test_update"
			assert.NoError(t, repo.UpdateClass(&class))
			assert.True(t, updated.Before(class.Updated))

			check, err := repo.GetClass(class.Id)
			assert.NoError(t, err)
			assert.Equal(t, class.Slug, check.Slug)
		})

		t.Run("Delete", func(t *testing.T) {
			class := model.Class{}
			assert.NoError(t, repo.InsertClass(&class))
			assert.NoError(t, repo.DeleteClass(class.Id))
			_, err := repo.GetClass(class.Id)
			// Make sure a "no documents in result" error pops out
			assert.Error(t, err)
		})
	})

	t.Run("Document", func(t *testing.T) {
		// Class used for tests below
		class := model.Class{
			Slug: "class_test",
		}
		assert.NoError(t, repo.InsertClass(&class))

		t.Run("Create", func(t *testing.T) {
			doc := model.Document{
				Slug: "create_test",
			}
			assert.NoError(t, repo.InsertDocument(&doc))
			assert.False(t, doc.Id.IsZero())
			assert.False(t, doc.Created.IsZero())
			assert.False(t, doc.Updated.IsZero())
		})

		t.Run("Read", func(t *testing.T) {
			doc := model.Document{
				ClassId: class.Id,
				Slug:    "read_test",
			}
			assert.NoError(t, repo.InsertDocument(&doc))

			check, err := repo.GetDocument(doc.Id)
			assert.NoError(t, err)
			assert.Equal(t, doc.ClassId, check.ClassId)
			assert.Equal(t, check.Class.Id, class.Id)
		})

		t.Run("Update", func(t *testing.T) {
			doc := model.Document{
				ClassId: class.Id,
				Slug:    "update_test",
			}
			assert.NoError(t, repo.InsertDocument(&doc))

			updated := doc.Updated
			doc.Slug = "update_test_update"
			assert.NoError(t, repo.UpdateDocument(&doc))
			assert.True(t, updated.Before(doc.Updated))

			check, err := repo.GetDocument(doc.Id)
			assert.NoError(t, err)
			assert.Equal(t, doc.Slug, check.Slug)
		})

		t.Run("Delete", func(t *testing.T) {
			doc := model.Document{}
			assert.NoError(t, repo.InsertDocument(&doc))
			assert.NoError(t, repo.DeleteDocument(doc.Id))

			_, err := repo.GetDocument(doc.Id)
			assert.Error(t, err)
		})
	})
}