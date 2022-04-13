package gocms

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ClassId   primitive.ObjectID `bson:"class_id"`
	ParentId  primitive.ObjectID `bson:"parent_id"`
	Title     string
	Slug      string
	Created   time.Time
	Updated   time.Time
	Published time.Time
	Values    map[string]interface{}
}

func (d Document) Value(key string) interface{} {
	switch key {
	case "id":
		return d.Id
	case "class_id":
		return d.ClassId
	case "parent_id":
		return d.ParentId
	case "title":
		return d.Title
	case "slug":
		return d.Slug
	case "published":
		return d.Published
	default:
		if v, ok := d.Values[key]; ok {
			return v
		}
	}
	return nil
}

type DocumentList struct {
	Total     int64
	Documents []Document
}

type DocumentListParams struct {
	ClassId primitive.ObjectID
	Page    int64
	Size    int64
}

type DocumentRepository interface {
	DeleteDocument(primitive.ObjectID) error
	GetChildDocumentBySlug(primitive.ObjectID, string) (Document, error)
	GetClassDocumentBySlug(primitive.ObjectID, string) (Document, error)
	GetDocumentList(DocumentListParams) (DocumentList, error)
	GetDocumentById(primitive.ObjectID) (Document, error)
	InsertDocument(*Document) error
	UpdateDocument(*Document) error
}

type DocumentService interface {
	Delete(Document) error
	GetById(primitive.ObjectID) (Document, error)
	GetChildBySlug(primitive.ObjectID, string) (Document, error)
	GetClassChildBySlug(primitive.ObjectID, string) (Document, error)
	Insert(*Document) error
	List(DocumentListParams) (DocumentList, error)
	Update(*Document) error
}

type documentService struct {
	repo DocumentRepository
}

func (p DocumentListParams) Offset() (offset int64) {
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Size
	}
	return
}

func NewDocumentService(repo DocumentRepository) DocumentService {
	return documentService{
		repo: repo,
	}
}

func (s documentService) Delete(doc Document) error {
	return s.repo.DeleteDocument(doc.Id)
}

func (s documentService) GetById(id primitive.ObjectID) (Document, error) {
	return s.repo.GetDocumentById(id)
}

func (s documentService) GetChildBySlug(parentId primitive.ObjectID, slug string) (Document, error) {
	return s.repo.GetChildDocumentBySlug(parentId, slug)
}

func (s documentService) GetClassChildBySlug(classId primitive.ObjectID, slug string) (Document, error) {
	return s.repo.GetClassDocumentBySlug(classId, slug)
}

func (s documentService) Insert(doc *Document) error {
	if err := s.Validate(doc); err != nil {
		return err
	}

	if !doc.Id.IsZero() {
		return fmt.Errorf("document already has an ID")
	}

	if doc.ParentId.IsZero() {
		check, err := s.GetClassChildBySlug(doc.ClassId, doc.Slug)
		if err == nil {
			return fmt.Errorf("slug %s already exists in %s", doc.Slug, check.Id.Hex())
		}
	}

	if !doc.ParentId.IsZero() {
		check, err := s.GetChildBySlug(doc.ParentId, doc.Slug)
		if err == nil {
			return fmt.Errorf("slug %s already exists in %s", doc.Slug, check.Id.Hex())
		}
	}

	return s.repo.InsertDocument(doc)
}

func (s documentService) List(params DocumentListParams) (DocumentList, error) {
	return s.repo.GetDocumentList(params)
}

func (s documentService) Update(doc *Document) error {
	if err := s.Validate(doc); err != nil {
		return err
	}

	if doc.Id.IsZero() {
		return fmt.Errorf("document has no ID")
	}

	if doc.ParentId.IsZero() {
		check, err := s.GetClassChildBySlug(doc.ClassId, doc.Slug)
		if err == nil && check.Id != doc.Id {
			return fmt.Errorf("slug %s already exists in %s", doc.Slug, check.Id.Hex())
		}
	}

	if !doc.ParentId.IsZero() {
		check, err := s.GetChildBySlug(doc.ParentId, doc.Slug)
		if err == nil {
			return fmt.Errorf("slug %s already exists in %s", doc.Slug, check.Id.Hex())
		}
	}

	return s.repo.UpdateDocument(doc)
}

func (s documentService) Validate(doc *Document) (err error) {
	if doc.ClassId.IsZero() {
		return fmt.Errorf("document requires a class ID")
	}

	if doc.Slug == "" {
		return fmt.Errorf("document requires a slug")
	}

	return
}
