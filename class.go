package gocms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `json:"name"`
	Slug    string             `json:"slug"`
	Created time.Time          `json:"created"`
	Updated time.Time          `json:"updated"`
	Fields  []Field            `json:"fields"`
}
