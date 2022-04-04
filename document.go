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
