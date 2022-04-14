package server

import (
	"fmt"
	"testing"

	"github.com/jbaikge/gocms"
	"github.com/zeebo/assert"
)

func TestEmptyTable(t *testing.T) {
	table := NewTable(gocms.Class{}, []gocms.Document{})
	assert.Equal(t, 1, len(table.Header()))
	assert.Equal(t, 0, len(table.Body()))
}

func TestPopulatedTable(t *testing.T) {
	class := gocms.Class{
		TableLabels: "Title A B C",
		TableFields: "title a b c",
	}
	docs := make([]gocms.Document, 3)
	for i := range docs {
		docs[i] = gocms.Document{
			Title: fmt.Sprintf("%d.0", i),
			Values: map[string]interface{}{
				"a": fmt.Sprintf("%d.1", i),
				"b": fmt.Sprintf("%d.2", i),
				"c": fmt.Sprintf("%d.3", i),
				"d": fmt.Sprintf("%d.4", i),
			},
		}
	}
	table := NewTable(class, docs)

	headers := table.Header()
	assert.Equal(t, 4, len(headers))
	for i, exp := range []string{"Title", "A", "B", "C"} {
		assert.Equal(t, exp, headers[i])
	}

	body := table.Body()
	assert.Equal(t, 3, len(body))
	for r, row := range body {
		assert.Equal(t, 4, len(row.Columns))
		for c, value := range row.Columns {
			assert.Equal(t, fmt.Sprintf("%d.%d", r, c), value)
		}
	}
}
