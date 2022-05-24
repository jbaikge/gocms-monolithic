package repository

import (
	"github.com/jbaikge/gocms/models/class"
	"github.com/jbaikge/gocms/models/document"
	"github.com/jbaikge/gocms/models/user"
)

type Repository interface {
	class.ClassRepository
	document.DocumentRepository
	user.UserRepository

	// Only used for testing
	empty() error
}
