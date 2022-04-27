package gocms

import (
	"fmt"
	"testing"
	"time"

	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDocumentValue(t *testing.T) {
	doc := Document{
		Id:        primitive.NewObjectID(),
		ClassId:   primitive.NewObjectID(),
		ParentId:  primitive.NewObjectID(),
		Slug:      "my_document",
		Title:     "My Title",
		Published: time.Now(),
		Values: map[string]interface{}{
			"a": "value a",
			"b": "value b",
		},
	}
	assert.Equal(t, doc.Title, doc.Value("title"))
	assert.Equal(t, doc.Id, doc.Value("id"))
	assert.Equal(t, doc.ClassId, doc.Value("class_id"))
	assert.Equal(t, doc.ParentId, doc.Value("parent_id"))
	assert.Equal(t, doc.Slug, doc.Value("slug"))
	assert.Equal(t, doc.Published, doc.Value("published"))
	assert.Equal(t, doc.Values["a"], doc.Value("a"))
	assert.Equal(t, doc.Values["b"], doc.Value("b"))
	assert.Equal(t, nil, doc.Value("nil"))
}

var _ DocumentRepository = mockDocumentRepository{}

type mockDocumentRepository struct {
	byId         map[primitive.ObjectID]Document
	byClassId    map[primitive.ObjectID][]Document
	byClassSlug  map[string]Document
	byParentSlug map[string]Document
}

func NewMockDocumentRepository() mockDocumentRepository {
	return mockDocumentRepository{
		byId:         make(map[primitive.ObjectID]Document),
		byClassId:    make(map[primitive.ObjectID][]Document),
		byClassSlug:  make(map[string]Document),
		byParentSlug: make(map[string]Document),
	}
}

func (r mockDocumentRepository) DeleteDocument(id primitive.ObjectID) (err error) {
	doc, ok := r.byId[id]
	if !ok {
		// Silent failure
		return
	}
	delete(r.byId, id)
	delete(r.byClassSlug, r.slugKey(doc.ClassId, doc.Slug))
	delete(r.byParentSlug, r.slugKey(doc.ParentId, doc.Slug))
	return
}

func (r mockDocumentRepository) GetDocumentById(id primitive.ObjectID) (doc Document, err error) {
	doc, ok := r.byId[id]
	if !ok {
		err = fmt.Errorf("document not found: %s", id)
	}
	return
}

func (r mockDocumentRepository) GetChildDocumentBySlug(parentId primitive.ObjectID, slug string) (doc Document, err error) {
	doc, ok := r.byParentSlug[r.slugKey(parentId, slug)]
	if ok {
		return
	}

	err = fmt.Errorf("document not found: %s", r.slugKey(parentId, slug))
	return
}

func (r mockDocumentRepository) GetClassDocumentBySlug(classId primitive.ObjectID, slug string) (doc Document, err error) {
	doc, ok := r.byClassSlug[r.slugKey(classId, slug)]
	if ok {
		return
	}

	err = fmt.Errorf("document not found: %s", r.slugKey(classId, slug))
	return
}

func (r mockDocumentRepository) GetDocumentList(params DocumentListParams) (list DocumentList, err error) {
	docs, ok := r.byClassId[params.ClassId]
	if !ok {
		err = fmt.Errorf("no documents for class: %s", params.ClassId.Hex())
	}

	offset := params.Offset()
	if offset >= int64(len(docs)) {
		err = fmt.Errorf("offset out of bounds: %d (length: %d)", offset, len(docs))
	}

	end := offset + params.Size
	if end > int64(len(docs)) {
		end = int64(len(docs))
	}

	list.Total = int64(len(docs))
	list.Documents = docs[offset:end]

	return
}

func (r mockDocumentRepository) InsertDocument(doc *Document) (err error) {
	doc.Id = primitive.NewObjectID()
	r.byId[doc.Id] = *doc
	r.byClassSlug[r.slugKey(doc.ClassId, doc.Slug)] = *doc
	r.byParentSlug[r.slugKey(doc.ParentId, doc.Slug)] = *doc
	if _, ok := r.byClassId[doc.ClassId]; !ok {
		r.byClassId[doc.ClassId] = make([]Document, 0, 16)
	}
	r.byClassId[doc.ClassId] = append(r.byClassId[doc.ClassId], *doc)
	return
}

func (r mockDocumentRepository) UpdateDocument(doc *Document) (err error) {
	if err = r.DeleteDocument(doc.Id); err != nil {
		return
	}
	r.byId[doc.Id] = *doc
	r.byClassSlug[r.slugKey(doc.ClassId, doc.Slug)] = *doc
	r.byParentSlug[r.slugKey(doc.ParentId, doc.Slug)] = *doc
	return
}

func (r mockDocumentRepository) slugKey(id primitive.ObjectID, slug string) string {
	return id.Hex() + "_" + slug
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

	t.Run("GetBySlug", func(t *testing.T) {
		service := NewDocumentService(NewMockDocumentRepository())

		doc := Document{
			ClassId:  primitive.NewObjectID(),
			ParentId: primitive.NewObjectID(),
			Slug:     "test",
		}
		assert.NoError(t, service.Insert(&doc))

		t.Run("Class ID", func(t *testing.T) {
			byClass, err := service.GetClassChildBySlug(doc.ClassId, doc.Slug)
			assert.NoError(t, err)
			assert.Equal(t, doc.Id, byClass.Id)
		})

		t.Run("Parent ID", func(t *testing.T) {
			byParent, err := service.GetChildBySlug(doc.ParentId, doc.Slug)
			assert.NoError(t, err)
			assert.Equal(t, doc.Id, byParent.Id)
		})
	})

	t.Run("Insert", func(t *testing.T) {
		service := NewDocumentService(NewMockDocumentRepository())
		classId := primitive.NewObjectID()
		parentId := primitive.NewObjectID()

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
				Document{
					Id:      primitive.NewObjectID(),
					ClassId: primitive.NewObjectID(),
					Slug:    "test",
				},
			},
			{
				"Parent ID & Slug",
				false,
				Document{
					ClassId:  primitive.NewObjectID(),
					ParentId: parentId,
					Slug:     "test",
				},
			},
			{
				"Parent ID & Dupe Slug",
				true,
				Document{
					ClassId:  primitive.NewObjectID(),
					ParentId: parentId,
					Slug:     "test",
				},
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

		t.Run("Class Slug Takeover", func(t *testing.T) {
			banana.Slug = orange.Slug
			defer func() {
				banana.Slug = "banana"
			}()

			assert.Error(t, service.Update(&banana))
		})

		t.Run("New Class, Dupe Slug", func(t *testing.T) {
			banana.ClassId = primitive.NewObjectID()
			banana.Slug = orange.Slug
			defer func() {
				banana.ClassId = classId
				banana.Slug = "banana"
			}()

			assert.NoError(t, service.Update(&banana))
		})

		t.Run("Same Parent, Slug Frob", func(t *testing.T) {
			banana.ParentId = primitive.NewObjectID()
			orange.ParentId = banana.ParentId
			defer func() {
				banana.ParentId = primitive.NilObjectID
				orange.ParentId = banana.ParentId
				orange.Slug = "orange"
			}()

			assert.NoError(t, service.Update(&banana))
			assert.NoError(t, service.Update(&orange))

			orange.Slug = banana.Slug
			assert.Error(t, service.Update(&orange))
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

	t.Run("List", func(t *testing.T) {
		service := NewDocumentService(NewMockDocumentRepository())

		classId := primitive.NewObjectID()
		ids := make([]primitive.ObjectID, 3)

		for i := range ids {
			doc := Document{
				ClassId: classId,
				Slug:    fmt.Sprintf("test_%d", i),
			}
			assert.NoError(t, service.Insert(&doc))
			ids[i] = doc.Id
		}

		params := DocumentListParams{
			ClassId: classId,
			Size:    2,
			Page:    1,
		}
		page1, err := service.List(params)
		assert.NoError(t, err)
		assert.Equal(t, 3, page1.Total)
		assert.Equal(t, 2, len(page1.Documents))
		for i := range ids[0:2] {
			assert.Equal(t, ids[i], page1.Documents[i].Id)
		}

		params.Page = 2
		page2, err := service.List(params)
		assert.NoError(t, err)
		assert.Equal(t, 3, page2.Total)
		assert.Equal(t, 1, len(page2.Documents))
		for i := range ids[2:3] {
			assert.Equal(t, ids[i], page1.Documents[i].Id)
		}
	})
}
