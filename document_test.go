package gocms

import (
	"fmt"
	"testing"

	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ DocumentRepository = mockDocumentRepository{}

type mockDocumentRepository struct {
	byId map[primitive.ObjectID]Document
}

func NewMockDocumentRepository() mockDocumentRepository {
	return mockDocumentRepository{
		byId: make(map[primitive.ObjectID]Document),
	}
}

func (r mockDocumentRepository) DeleteDocument(id primitive.ObjectID) (err error) {
	_, ok := r.byId[id]
	if !ok {
		// Silent failure
		return
	}
	delete(r.byId, id)
	return
}

func (r mockDocumentRepository) GetDocumentById(id primitive.ObjectID) (doc Document, err error) {
	doc, ok := r.byId[id]
	if !ok {
		err = fmt.Errorf("document not found: %s", id)
	}
	return
}

func (r mockDocumentRepository) InsertDocument(doc *Document) (err error) {
	doc.Id = primitive.NewObjectID()
	r.byId[doc.Id] = *doc
	return
}

func (r mockDocumentRepository) UpdateDocument(doc *Document) (err error) {
	if err = r.DeleteDocument(doc.Id); err != nil {
		return
	}
	r.byId[doc.Id] = *doc
	return
}

func TestDocumentService(t *testing.T) {

	t.Run("GetById", func(t *testing.T) {
		service := NewDocumentService(NewMockDocumentRepository())

		doc := Document{ClassId: primitive.NewObjectID(), Slug: "test"}
		assert.NoError(t, service.Insert(&doc))

		check, err := service.GetById(doc.Id)
		assert.NoError(t, err)
		assert.Equal(t, doc.Id, check.Id)
	})

	t.Run("Insert", func(t *testing.T) {
		service := NewDocumentService(NewMockDocumentRepository())
		classId := primitive.NewObjectID()

		tests := []struct {
			Name     string
			Error    bool
			Document Document
		}{
			{
				"Slug Only",
				true,
				Document{Slug: "test"},
			},
			{
				"Class ID Only",
				true,
				Document{ClassId: primitive.NewObjectID()},
			},
			{
				"Class ID & Slug",
				false,
				Document{ClassId: classId, Slug: "test"},
			},
			{
				"Duplicate Slug",
				true,
				Document{ClassId: classId, Slug: "test"},
			},
			{
				"Same Class, New Slug",
				false,
				Document{ClassId: classId, Slug: "new_test"},
			},
			{
				"Same Slug, New Class",
				false,
				Document{ClassId: primitive.NewObjectID(), Slug: "test"},
			},
			{
				"Existing ID",
				true,
				Document{Id: primitive.NewObjectID(), ClassId: primitive.NewObjectID(), Slug: "test"},
			},
		}

		for _, test := range tests {
			t.Run(test.Name, func(t *testing.T) {
				err := service.Insert(&test.Document)
				if test.Error {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		service := NewDocumentService(NewMockDocumentRepository())

		t.Run("No ID", func(t *testing.T) {
			doc := Document{ClassId: primitive.NewObjectID(), Slug: "test"}
			assert.Error(t, service.Update(&doc))
		})

		classId := primitive.NewObjectID()

		banana := Document{ClassId: classId, Slug: "banana"}
		assert.NoError(t, service.Insert(&banana))

		orange := Document{ClassId: classId, Slug: "orange"}
		assert.NoError(t, service.Insert(&orange))

		t.Run("Blank Slug", func(t *testing.T) {
			banana.Slug = ""
			defer func() {
				banana.Slug = "banana"
			}()

			assert.Error(t, service.Update(&banana))
		})

		t.Run("Slug Takeover", func(t *testing.T) {
			banana.Slug = "orange"
			defer func() {
				banana.Slug = "banana"
			}()

			assert.Error(t, service.Update(&banana))
		})

		t.Run("New Class, Dupe Slug", func(t *testing.T) {
			banana.ClassId = primitive.NewObjectID()
			banana.Slug = "orange"
			defer func() {
				banana.ClassId = classId
				banana.Slug = "banana"
			}()

			assert.NoError(t, service.Update(&banana))
		})
	})

	t.Run("Delete", func(t *testing.T) {
		service := NewDocumentService(NewMockDocumentRepository())

		doc := Document{ClassId: primitive.NewObjectID(), Slug: "test"}
		assert.NoError(t, service.Insert(&doc))
		assert.NoError(t, service.Delete(doc))
		// Do it once more to make sure it fails silently
		assert.NoError(t, service.Delete(doc))

		// Make sure document no longer exists
		_, err := service.GetById(doc.Id)
		assert.Error(t, err)
	})
}
