package server

import (
	"strings"

	"github.com/jbaikge/gocms"
)

// Manages data to build HTML tables
type Table struct {
	class     *gocms.Class
	documents []gocms.Document
}

type TableRow struct {
	Document gocms.Document
	Columns  []string
}

func NewTable(class *gocms.Class, docs []gocms.Document) Table {
	return Table{
		class:     class,
		documents: docs,
	}
}

func (t Table) Header() []string {
	return strings.Fields(t.class.TableLabels)
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
			raw := doc.Value(name)
			if raw == nil {
				// Silently return a blank string
				continue
			}

			rows[i].Columns[n] = t.class.Field(name).Apply(raw)
		}
	}

	return rows
}
