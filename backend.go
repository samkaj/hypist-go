package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func main() {
  setupDb()
  startServer()
}

func setupDb() {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	dbClient, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect(context.Background())
}

func startServer() {
	r := gin.Default()
	r.POST("/users", PostUser)
	r.Run()
}
