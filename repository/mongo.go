package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jbaikge/gocms"
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
	filter := bson.M{"_id": id}
	_, err = m.classes.DeleteOne(m.context, filter)
	if err != nil {
		return
	}
	// Can include a check for the result and DeletedCount later if it is useful
	return
}

func (m mongoRepository) GetAllClasses() (classes []gocms.Class, err error) {
	cursor, err := m.classes.Find(m.context, bson.M{})
	if err != nil {
		return
	}
	err = cursor.All(m.context, &classes)
	return
}

func (m mongoRepository) GetClassById(id primitive.ObjectID) (class gocms.Class, err error) {
	filter := bson.M{"_id": id}
	err = m.classes.FindOne(m.context, filter).Decode(&class)
	return
}

func (m mongoRepository) GetClassBySlug(slug string) (class gocms.Class, err error) {
	filter := bson.M{"slug": slug}
	err = m.classes.FindOne(m.context, filter).Decode(&class)
	return
}

func (m mongoRepository) InsertClass(class *gocms.Class) (err error) {
	now := time.Now()
	class.Created = now
	class.Updated = now
	result, err := m.classes.InsertOne(m.context, class)
	if err != nil {
		return
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("Unable to cast newly inserted Class ID to ObjectID")
	}
	class.Id = id
	return
}

func (m mongoRepository) UpdateClass(class *gocms.Class) (err error) {
	class.Updated = time.Now()
	filter := bson.M{"_id": class.Id}
	result, err := m.classes.ReplaceOne(m.context, filter, class)
	if err != nil {
		return
	}
	if result.MatchedCount == 0 {
		return errors.New("Did not match a Class to update")
	}
	return
}

func (m mongoRepository) DeleteDocument(id primitive.ObjectID) (err error) {
	filter := bson.M{"_id": id}
	_, err = m.documents.DeleteOne(m.context, filter)
	if err != nil {
		return
	}
	return
}

func (m mongoRepository) GetDocumentById(id primitive.ObjectID) (doc gocms.Document, err error) {
	filter := bson.M{"_id": id}
	err = m.documents.FindOne(m.context, filter).Decode(&doc)
	return
}

func (m mongoRepository) InsertDocument(doc *gocms.Document) (err error) {
	now := time.Now()
	doc.Created = now
	doc.Updated = now

	result, err := m.documents.InsertOne(m.context, doc)
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("Unable to cast newly inserted Document ID to ObjectID")
	}

	doc.Id = id
	return
}

func (m mongoRepository) UpdateDocument(doc *gocms.Document) (err error) {
	doc.Updated = time.Now()

	filter := bson.M{"_id": doc.Id}
	result, err := m.documents.ReplaceOne(m.context, filter, doc)
	if err != nil {
		return
	}
	if result.MatchedCount == 0 {
		return errors.New("Did not match a Document to update")
	}
	return
}
