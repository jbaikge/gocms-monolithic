package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms/models/class"
	"github.com/jbaikge/gocms/models/document"
	"github.com/jbaikge/gocms/repository"
	"github.com/jbaikge/gocms/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	dbHost := "localhost:27017"
	if dbHostEnv := os.Getenv("DB_HOST"); dbHostEnv != "" {
		dbHost = dbHostEnv
	}

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+dbHost))
	if err != nil {
		log.Fatalf("Unable to create client %v", err)
	}

	db := client.Database("gocms-web")

	repo := repository.NewMongo(ctx, db)
	classService := class.NewClassService(repo)
	documentService := document.NewDocumentService(repo)

	router := gin.Default()
	router.SetTrustedProxies(nil)
	s := server.New(router, classService, documentService)
	panic(s.Run(":8080"))
}
