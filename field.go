package gocms

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
