package gocms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ClassId   primitive.ObjectID `bson:"class_id"`
	Slug      string
	Created   time.Time
	Updated   time.Time
	Published time.Time
	Class     Class `bson:"-"`
	Values    map[string]interface{}
}
