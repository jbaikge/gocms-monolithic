package gocms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ClassId   primitive.ObjectID `bson:"class_id"`
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

func (s documentService) Delete(doc Document) error {
	return s.repo.DeleteDocument(doc.Id)
}

func (s documentService) GetById(id primitive.ObjectID) (Document, error) {
	return s.repo.GetDocumentById(id)
}

func (s documentService) Insert(doc *Document) error {
	return s.repo.InsertDocument(doc)
}

func (s documentService) Update(doc *Document) error {
	return s.repo.UpdateDocument(doc)
}
