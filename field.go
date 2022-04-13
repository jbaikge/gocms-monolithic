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
	Type            string             `json:"type"`
	Name            string             `json:"name"`
	Label           string             `json:"label"`
	Min             string             `json:"min" bson:",omitempty"`
	Max             string             `json:"max" bson:",omitempty"`
	Step            string             `json:"step" bson:",omitempty"`
	Format          string             `json:"format" bson:",omitempty"`
	Options         string             `json:"options" bson:",omitempty"`
	DataSourceId    primitive.ObjectID `json:"data_source_id" bson:"data_source_id,omitempty" form:"data_source_id"`
	DataSourceValue string             `json:"data_source_value" bson:"data_source_value,omitempty" form:"data_source_value"`
	DataSourceLabel string             `json:"data_source_label" bson:"data_source_label,omitempty" form:"data_source_label"`
}

// Takes in any value from a Document.Values item and converts it based on the
// field type, then optionally formats the value if defined
func (f Field) Apply(value interface{}) string {
	if s, ok := value.(string); ok {
		if f.Format == "" {
			return s
		}
		if f.Type == TypeDate {
			t, _ := time.Parse("2006-01-02", s)
			return t.Format(f.Format)
		}
		if f.Type == TypeDateTime {
			t, _ := time.Parse("2006-01-02T15:04", s)
			return t.Format(f.Format)
		}
		if f.Type == TypeTime {
			t, _ := time.Parse("15:04", s)
			return t.Format(f.Format)
		}
	}
	return "-nil-"
}
