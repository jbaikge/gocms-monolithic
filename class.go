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

type classService struct {
	repo ClassRepository
}

func NewClassService(repo ClassRepository) ClassService {
	return classService{
		repo: repo,
	}
}

func (s classService) Delete(class Class) error {
	return s.repo.DeleteClass(class.Id)
}

func (s classService) GetById(id primitive.ObjectID) (Class, error) {
	return s.repo.GetClassById(id)
}

func (s classService) GetBySlug(slug string) (Class, error) {
	return s.repo.GetClassBySlug(slug)
}

func (s classService) Insert(class *Class) error {
	return s.repo.InsertClass(class)
}

func (s classService) Update(class *Class) error {
	return s.repo.UpdateClass(class)
}
