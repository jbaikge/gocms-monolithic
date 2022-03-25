package document

type Option struct {
	Value    string
	Label    string
	Selected bool
}

type Type int

type Field struct {
	Id      int64
	Name    string
	Label   string
	Type    Type
	Options []Option
	Fields  []Field
}

type Class struct {
	Id     int64
	Slug   string
	Fields []Field
}

type Document struct {
	Id     int64
	Slug   string
	Class  Class
	Values map[string]interface{}
}

const (
	Text Type = iota
	TextArea
	Number
	Date
	Time
	Email
	Select
	MultiSelect
	Section
)
