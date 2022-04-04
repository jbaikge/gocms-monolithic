package gocms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `json:"name"`
	Slug    string             `json:"slug"`
	Created time.Time          `json:"created"`
	Updated time.Time          `json:"updated"`
	Fields  []Field            `json:"fields"`
}

type ClassRepository interface {
	DeleteClass(primitive.ObjectID) error
	GetClassById(primitive.ObjectID) (Class, error)
	GetClassBySlug(string) (Class, error)
	InsertClass(*Class) error
	UpdateClass(*Class) error
}

type ClassService interface {
	Delete(Class) error
	GetById(primitive.ObjectID) (Class, error)
	GetBySlug(string) (Class, error)
	Insert(*Class) error
	Update(*Class) error
}
