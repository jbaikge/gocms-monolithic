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
			var value string
			switch name {
			case "title":
				value = doc.Title
			case "slug":
				value = doc.Slug
			case "published":
				value = doc.Published.String() // TODO format published date
			default:
				val, ok := doc.Values[name]
				if !ok {
					// Silently return a blank string
					break
				}
				value = t.class.Field(name).Apply(val)
			}
			rows[i].Columns[n] = value
		}
	}

	return rows
}
