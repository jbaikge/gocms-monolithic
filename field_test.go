package gocms

import (
	"testing"
	"time"

	"github.com/zeebo/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestEmptyFieldApply(t *testing.T) {
	field := Field{}
	value := "value"
	assert.Equal(t, value, field.Apply(value))
}

func TestFieldApply(t *testing.T) {
	now := time.Now()
	objectId := primitive.NewObjectID()

	table := []struct {
		Name   string
		Type   string
		Format string
		Value  interface{}
		Expect string
	}{
		{"Date", TypeDate, "Jan 2, 2006", "2022-04-14", "Apr 14, 2022"},
		{"Date & Time", TypeDateTime, "Jan 2, 2006 3:04 pm", "2022-04-14T12:08", "Apr 14, 2022 12:08 pm"},
		{"Email", TypeEmail, "", "test@test.com", "test@test.com"},
		{"Multi-Select", TypeMultiSelect, "", []string{"a", "b"}, "-nil-"},
		{"Number String", TypeNumber, "", "42", "42"},
		{"Number Number", TypeNumber, "", 42, "42"},
		{"Select", TypeSelect, "", "option", "option"},
		{"Text", TypeText, "", "text", "text"},
		{"Textarea", TypeTextArea, "", "textarea", "textarea"},
		{"Time", TypeTime, "3:04 pm", "12:11", "12:11 pm"},
		{"TinyMCE", TypeTinyMCE, "", "tinymce", "tinymce"},
		{"time.Time", "", "", now, now.Format("Jan 2, 2006 3:04pm")},
		{"ObjectID", "", "", objectId, objectId.Hex()},
	}

	for _, test := range table {
		t.Run(test.Name, func(t *testing.T) {
			f := Field{
				Type:   test.Type,
				Format: test.Format,
			}
			assert.Equal(t, test.Expect, f.Apply(test.Value))
		})
	}
}

func TestFieldOptionList(t *testing.T) {
	expect := []FieldOption{
		{
			Value: "a",
			Label: "A",
		},
		{
			Value: "b",
			Label: "B",
		},
		{
			Value: "C",
			Label: "C",
		},
	}

	field := Field{
		Options: "a|A\nb|B\nC\n",
	}
	options := field.OptionList()

	assert.Equal(t, len(expect), len(options))
	for i := range options {
		assert.DeepEqual(t, expect[i], options[i])
	}
}
