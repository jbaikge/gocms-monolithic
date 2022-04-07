package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jbaikge/gocms"
	"github.com/jbaikge/gocms/repository"
	"github.com/jbaikge/gocms/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Unable to create client %v", err)
	}

	db := client.Database("gocms-web")

	repo := repository.NewMongo(ctx, db)
	classService := gocms.NewClassService(repo)

	router := gin.Default()
	router.SetTrustedProxies(nil)
	s := server.New("./web", router, classService)
	panic(s.Run(":8080"))
}
