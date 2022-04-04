package gocms

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
	TypeSelect      = "select"
	TypeText        = "text"
	TypeTextArea    = "textarea"
	TypeTime        = "time"
	TypeTinyMCE     = "tinymce"
	TypeUpload      = "upload"
)

type Field struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	Min     string `json:"min"`
	Max     string `json:"max"`
	Step    string `json:"step"`
	Options string `json:"options"`
}

type Class struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `json:"name"`
	Slug    string             `json:"slug"`
	Created time.Time          `json:"created"`
	Updated time.Time          `json:"updated"`
	Fields  []Field            `json:"fields"`
}

type Document struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ClassId   primitive.ObjectID `bson:"class_id"`
	Slug      string
	Created   time.Time
	Updated   time.Time
	Published time.Time
	Class     Class `bson:"-"`
	Values    map[string]interface{}
}
