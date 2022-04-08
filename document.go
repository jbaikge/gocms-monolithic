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
	Slug      string
	Created   time.Time
	Updated   time.Time
	Published time.Time
	Values    map[string]interface{}
}

type DocumentRepository interface {
	DeleteDocument(primitive.ObjectID) error
	GetDocumentById(primitive.ObjectID) (Document, error)
	InsertDocument(*Document) error
	UpdateDocument(*Document) error
}

type DocumentService interface {
	Delete(Document) error
	GetById(primitive.ObjectID) (Document, error)
	Insert(*Document) error
	Update(*Document) error
}

type documentService struct {
	repo DocumentRepository
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

func (s documentService) Insert(doc *Document) error {
	if err := s.Validate(doc); err != nil {
		return err
	}

	if !doc.Id.IsZero() {
		return fmt.Errorf("document already has an ID")
	}

	return s.repo.InsertDocument(doc)
}

func (s documentService) Update(doc *Document) error {
	if err := s.Validate(doc); err != nil {
		return err
	}

	if doc.Id.IsZero() {
		return fmt.Errorf("document has no ID")
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
