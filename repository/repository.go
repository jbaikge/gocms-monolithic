package repository

import (
	"github.com/jbaikge/gocms"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	// Class CRUD
	DeleteClass(primitive.ObjectID) error
	GetClass(primitive.ObjectID) (gocms.Class, error)
	InsertClass(*gocms.Class) error
	UpdateClass(*gocms.Class) error

	// Document CRUD
	DeleteDocument(primitive.ObjectID) error
	GetDocument(primitive.ObjectID) (gocms.Document, error)
	InsertDocument(*gocms.Document) error
	UpdateDocument(*gocms.Document) error
}
