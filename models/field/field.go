package gocms

import (
	"fmt"
	"strings"
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

type FieldOption struct {
	Value string
	Label string
}

type Field struct {
	Type            string             `json:"type"`
	Name            string             `json:"name"`
	Label           string             `json:"label"`
	Min             string             `json:"min" bson:",omitempty"`
	Max             string             `json:"max" bson:",omitempty"`
	Step            string             `json:"step" bson:",omitempty"`
	Format          string             `json:"format" bson:",omitempty"`
	Options         string             `json:"options" bson:",omitempty"`
	DataSourceId    primitive.ObjectID `json:"data_source_id" bson:"data_source_id,omitempty"`
	DataSourceValue string             `json:"data_source_value" bson:"data_source_value,omitempty"`
	DataSourceLabel string             `json:"data_source_label" bson:"data_source_label,omitempty"`
}

// Takes in any value from a Document.Values item and converts it based on the
// field type, then optionally formats the value if defined
func (f Field) Apply(value interface{}) string {
	switch v := value.(type) {
	case int:
		return fmt.Sprint(v)
	case string:
		if f.Format == "" {
			return v
		}
		if f.Type == TypeDate {
			t, _ := time.Parse("2006-01-02", v)
			return t.Format(f.Format)
		}
		if f.Type == TypeDateTime {
			t, _ := time.Parse("2006-01-02T15:04", v)
			return t.Format(f.Format)
		}
		if f.Type == TypeTime {
			t, _ := time.Parse("15:04", v)
			return t.Format(f.Format)
		}
	case primitive.ObjectID:
		return v.Hex()
	case time.Time:
		return v.Format("Jan 2, 2006 3:04pm")
	}
	return "-nil-"
}

// Converts the options text to an array for use in HTML templates to
// generate select options
func (f Field) OptionList() (options []FieldOption) {
	rows := strings.Split(strings.TrimSpace(f.Options), "\n")
	options = make([]FieldOption, len(rows))
	for i, row := range rows {
		s := strings.SplitN(row, "|", 2)
		if len(s) == 1 {
			s = append(s, s[0])
		}
		options[i] = FieldOption{
			Value: strings.TrimSpace(s[0]),
			Label: strings.TrimSpace(s[1]),
		}
	}
	return
}
