package repository

import (
	"github.com/jbaikge/gocms/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	// Class CRUD
	DeleteClass(primitive.ObjectID) error
	GetClass(primitive.ObjectID) (model.Class, error)
	InsertClass(*model.Class) error
	UpdateClass(*model.Class) error

	// Document CRUD
	DeleteDocument(primitive.ObjectID) error
	GetDocument(primitive.ObjectID) (model.Document, error)
	InsertDocument(*model.Document) error
	UpdateDocument(*model.Document) error
}
