package server

import (
	"strings"

	"github.com/jbaikge/gocms"
)

// Manages data to build HTML tables
type Table struct {
	class     gocms.Class
	documents []gocms.Document
}

// Holds necessary information for table rows including the document and
// the formatted strings for each column.
type TableRow struct {
	Document gocms.Document
	Columns  []string
}

// Creates a new table
func NewTable(class gocms.Class, docs []gocms.Document) Table {
	return Table{
		class:     class,
		documents: docs,
	}
}

func (t Table) Header() []string {
	headings := strings.Fields(t.class.TableLabels)
	if len(headings) == 0 {
		headings = []string{"Title"}
	}
	return headings
}

func (t Table) Body() (rows []TableRow) {
	rows = make([]TableRow, len(t.documents))

	names := strings.Fields(t.class.TableFields)
	if len(names) == 0 {
		names = []string{"title"}
	}

	for i, doc := range t.documents {
		rows[i].Document = doc
		rows[i].Columns = make([]string, len(names))
		for n, name := range names {
			rows[i].Columns[n] = t.class.Field(name).Apply(doc.Value(name))
		}
	}

	return rows
}
