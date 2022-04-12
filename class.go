package gocms

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Classes define a type of Document
type Class struct {
	Id            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Parents       []primitive.ObjectID `json:"parents"`
	Name          string               `json:"name" bson:"name" form:"name"`
	SingularName  string               `json:"singular_name" bson:"singular_name" form:"singular_name"`
	MenuLabel     string               `json:"menu_label" bson:"menu_label" form:"menu_label"`
	AddItemLabel  string               `json:"add_item_label" bson:"add_item_label" form:"add_item_label"`
	NewItemLabel  string               `json:"new_item_label" bson:"new_item_label" form:"new_item_label"`
	EditItemLabel string               `json:"edit_item_label" bson:"edit_item_label" form:"edit_item_label"`
	Slug          string               `json:"slug" bson:"slug" form:"slug"`
	TableLabels   string               `json:"table_labels" bson:"table_labels" form:"table_labels"`
	TableFields   string               `json:"table_fields" bson:"table_fields" form:"table_fields"`
	Created       time.Time            `json:"created"`
	Updated       time.Time            `json:"updated"`
	Fields        []Field              `json:"fields"`
}

func (c Class) Labels() []string {
	return strings.Fields(c.TableLabels)
}

// Repositories manage data storage and retrieval
type ClassRepository interface {
	DeleteClass(primitive.ObjectID) error
	GetAllClasses() ([]Class, error)
	GetClassById(primitive.ObjectID) (Class, error)
	GetClassBySlug(string) (Class, error)
	InsertClass(*Class) error
	UpdateClass(*Class) error
}

// Services manage business rules while interacting with repositories
type ClassService interface {
	All() ([]Class, error)
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

func (s classService) All() ([]Class, error) {
	return s.repo.GetAllClasses()
}

func (s classService) Delete(class Class) error {
	return s.repo.DeleteClass(class.Id)
}

func (s classService) GetById(id primitive.ObjectID) (Class, error) {
	return s.repo.GetClassById(id)
}

func (s classService) GetBySlug(slug string) (class Class, err error) {
	return s.repo.GetClassBySlug(slug)
}

func (s classService) Insert(class *Class) (err error) {
	if err = s.Validate(class); err != nil {
		return
	}

	if !class.Id.IsZero() {
		return fmt.Errorf("class already has an ID")
	}

	if check, err := s.GetBySlug(class.Slug); err == nil {
		return fmt.Errorf("slug %s already exists in %s", class.Slug, check.Id.Hex())
	}

	return s.repo.InsertClass(class)
}

func (s classService) Update(class *Class) (err error) {
	if err = s.Validate(class); err != nil {
		return
	}

	if class.Id.IsZero() {
		return fmt.Errorf("class has no ID")
	}

	if check, err := s.GetBySlug(class.Slug); err == nil && check.Id != class.Id {
		return fmt.Errorf("slug %s already exists in %s", class.Slug, check.Id.Hex())
	}

	return s.repo.UpdateClass(class)
}

func (s classService) Validate(class *Class) (err error) {
	if class.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if class.Slug == "" {
		return fmt.Errorf("slug is empty")
	}

	return
}
