package gocms

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	Type            string             `json:"type"`
	Name            string             `json:"name"`
	Label           string             `json:"label"`
	Min             string             `json:"min" bson:",omitempty"`
	Max             string             `json:"max" bson:",omitempty"`
	Step            string             `json:"step" bson:",omitempty"`
	Options         string             `json:"options" bson:",omitempty"`
	DataSourceId    primitive.ObjectID `json:"data_source_id" bson:"data_source_id,omitempty"`
	DataSourceValue string             `json:"data_source_value" bson:"data_source_value,omitempty"`
	DataSourceTitle string             `json:"data_source_title" bson:"data_source_title,omitempty"`
}
