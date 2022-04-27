package repository

import (
	"github.com/jbaikge/gocms/models/class"
	"github.com/jbaikge/gocms/models/document"
)

type Repository interface {
	class.ClassRepository
	document.DocumentRepository

	// Only used for testing
	empty() error
}
