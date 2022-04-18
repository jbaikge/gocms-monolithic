package repository

import (
	"github.com/jbaikge/gocms"
)

type Repository interface {
	gocms.ClassRepository
	gocms.DocumentRepository

	// Only used for testing
	empty() error
}
