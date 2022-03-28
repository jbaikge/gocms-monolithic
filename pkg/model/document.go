package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TypeDate        = "date"
	TypeDateTime    = "datetime"
	TypeEmail       = "email"
	TypeMultiSelect = "multiselect"
	TypeNumber      = "number"
	TypeSection     = "section"
	TypeSelect      = "select"
	TypeText        = "text"
	TypeTextArea    = "textarea"
	TypeTime        = "time"
)

type Option struct {
	Value    string
	Label    string
	Selected bool
}

type Field struct {
	Name    string
	Label   string
	Type    string
	Options []Option
	Fields  []Field
}

type Class struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Slug    string
	Created time.Time
	Updated time.Time
	Fields  []Field
}

type Document struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ClassId   primitive.ObjectID `bson:"class_id"`
	Slug      string
	Created   time.Time
	Updated   time.Time
	Published time.Time
	Class     *Class `bson:",omitempty"`
	Values    map[string]interface{}
}
