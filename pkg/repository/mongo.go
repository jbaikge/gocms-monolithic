package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jbaikge/gocms/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	context   context.Context
	db        *mongo.Database
	classes   *mongo.Collection
	documents *mongo.Collection
}

func NewMongo(ctx context.Context, db *mongo.Database) Repository {
	return &mongoRepository{
		context:   ctx,
		db:        db,
		classes:   db.Collection("classes"),
		documents: db.Collection("documents"),
	}
}

func (m mongoRepository) DeleteClass(id primitive.ObjectID) (err error) {
	return
}

func (m mongoRepository) GetClass(id primitive.ObjectID) (class model.Class, err error) {
	filter := bson.M{"_id": id}
	result := m.classes.FindOne(m.context, filter)
	err = result.Decode(&class)
	return
}

func (m mongoRepository) InsertClass(class *model.Class) (err error) {
	now := time.Now()
	class.Created = now
	class.Updated = now
	result, err := m.classes.InsertOne(m.context, class)
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("Unable to cast newly inserted ID to ObjectID")
	}
	class.Id = id
	return
}

func (m mongoRepository) UpdateClass(class *model.Class) (err error) {
	class.Updated = time.Now()
	filter := bson.M{"_id": class.Id}
	result, err := m.classes.ReplaceOne(m.context, filter, class)
	if err != nil {
		return
	}
	if result.MatchedCount == 0 {
		return errors.New("Did not match a document to update")
	}
	return
}
