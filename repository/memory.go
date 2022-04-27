package repository

import (
	"fmt"
	"sort"
	"time"

	"github.com/jbaikge/gocms/models/class"
	"github.com/jbaikge/gocms/models/document"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type sortClasses []class.Class

func (s sortClasses) Len() int           { return len(s) }
func (s sortClasses) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s sortClasses) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type memoryRepository struct {
	classes   []class.Class
	documents []document.Document
}

func NewMemory() Repository {
	return &memoryRepository{
		classes:   make([]class.Class, 0, 128),
		documents: make([]document.Document, 0, 128),
	}
}

func (r *memoryRepository) DeleteClass(id primitive.ObjectID) (err error) {
	for i, class := range r.classes {
		if class.Id == id {
			r.classes = append(r.classes[:i], r.classes[i+1:]...)
			break
		}
	}
	return
}

func (r *memoryRepository) GetAllClasses() ([]class.Class, error) {
	sort.Sort(sortClasses(r.classes))
	return r.classes, nil
}

func (r *memoryRepository) GetClassById(id primitive.ObjectID) (class class.Class, err error) {
	for _, c := range r.classes {
		if c.Id == id {
			return c, nil
		}
	}
	err = fmt.Errorf("class not found: %s", id.Hex())
	return
}

func (r *memoryRepository) GetClassBySlug(slug string) (class class.Class, err error) {
	for _, c := range r.classes {
		if c.Slug == slug {
			return c, nil
		}
	}
	err = fmt.Errorf("class not found: %s", slug)
	return
}

func (r *memoryRepository) InsertClass(class *class.Class) (err error) {
	class.Id = primitive.NewObjectID()
	now := time.Now()
	class.Created = now
	class.Updated = now
	r.classes = append(r.classes, *class)
	return
}

func (r *memoryRepository) UpdateClass(class *class.Class) (err error) {
	for i, c := range r.classes {
		if c.Id == class.Id {
			class.Updated = time.Now()
			r.classes[i] = *class
			return
		}
	}
	return fmt.Errorf("class not found: %s", class.Id.Hex())
}

func (r *memoryRepository) DeleteDocument(id primitive.ObjectID) (err error) {
	for i, doc := range r.documents {
		if doc.Id == id {
			r.documents = append(r.documents[:i], r.documents[i+1:]...)
			break
		}
	}
	return
}

func (r *memoryRepository) GetChildDocumentBySlug(parentId primitive.ObjectID, slug string) (doc document.Document, err error) {
	for _, d := range r.documents {
		if d.ParentId == parentId && d.Slug == slug {
			return d, nil
		}
	}
	err = fmt.Errorf("document not found for %s-%s", parentId.Hex(), slug)
	return
}

func (r *memoryRepository) GetClassDocumentBySlug(classId primitive.ObjectID, slug string) (doc document.Document, err error) {
	for _, d := range r.documents {
		if d.ClassId == classId && d.Slug == slug {
			return d, nil
		}
	}
	err = fmt.Errorf("document not found for %s-%s", classId.Hex(), slug)
	return
}

func (r *memoryRepository) GetDocumentList(params document.DocumentListParams) (list document.DocumentList, err error) {
	docs := make([]document.Document, 0, len(r.documents))
	for _, doc := range r.documents {
		if doc.ClassId == params.ClassId {
			docs = append(docs, doc)
		}
	}

	list.Total = int64(len(docs))
	if list.Total == 0 {
		return
	}

	end := params.Offset() + params.Size
	if end > list.Total {
		end = list.Total
	}
	list.Documents = docs[params.Offset():end]

	return
}

func (r *memoryRepository) GetDocumentById(id primitive.ObjectID) (doc document.Document, err error) {
	for _, d := range r.documents {
		if d.Id == id {
			return d, nil
		}
	}
	err = fmt.Errorf("document not found: %s", id.Hex())
	return
}

func (r *memoryRepository) InsertDocument(doc *document.Document) (err error) {
	doc.Id = primitive.NewObjectID()
	now := time.Now()
	doc.Created = now
	doc.Updated = now
	r.documents = append(r.documents, *doc)
	return
}

func (r *memoryRepository) UpdateDocument(doc *document.Document) (err error) {
	for i, d := range r.documents {
		if d.Id == doc.Id {
			doc.Updated = time.Now()
			r.documents[i] = *doc
			return
		}
	}
	return fmt.Errorf("document not found: %s", doc.Id.Hex())
}

func (r *memoryRepository) empty() (err error) {
	r.classes = r.classes[:0]
	r.documents = r.documents[:0]
	return
}
