package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jbaikge/gocms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	// Can include a check for the result and DeletedCount later if it is useful
	return
}

func (m mongoRepository) GetAllClasses() (classes []gocms.Class, err error) {
	filter := bson.D{}
	sort := bson.D{bson.E{Key: "name", Value: 1}}
	opts := options.Find().SetSort(sort)

	cursor, err := m.classes.Find(m.context, filter, opts)
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
	return
}

func (m mongoRepository) GetDocumentById(id primitive.ObjectID) (doc gocms.Document, err error) {
	filter := bson.M{"_id": id}
	err = m.documents.FindOne(m.context, filter).Decode(&doc)
	return
}

func (m mongoRepository) GetChildDocumentBySlug(id primitive.ObjectID, slug string) (doc gocms.Document, err error) {
	filter := bson.D{{Key: "parent_id", Value: id}, {Key: "slug", Value: slug}}
	err = m.documents.FindOne(m.context, filter).Decode(&doc)
	return
}

func (m mongoRepository) GetClassDocumentBySlug(id primitive.ObjectID, slug string) (doc gocms.Document, err error) {
	filter := bson.D{{Key: "class_id", Value: id}, {Key: "slug", Value: slug}}
	err = m.documents.FindOne(m.context, filter).Decode(&doc)
	return
}

func (m mongoRepository) GetDocumentList(params gocms.DocumentListParams) (list gocms.DocumentList, err error) {
	filter := bson.D{{Key: "class_id", Value: params.ClassId}}

	countOpts := options.Count()
	list.Total, err = m.documents.CountDocuments(m.context, filter, countOpts)
	if err != nil {
		return
	}
	if list.Total == 0 {
		return
	}

	findOpts := options.Find().SetLimit(params.Size).SetSkip(params.Offset())
	cursor, err := m.documents.Find(m.context, filter, findOpts)
	if err != nil {
		return
	}
	err = cursor.All(m.context, &list.Documents)
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

func (m mongoRepository) empty() (err error) {
	if err := m.documents.Drop(m.context); err != nil {
		return err
	}
	if err := m.classes.Drop(m.context); err != nil {
		return err
	}
	return
}
