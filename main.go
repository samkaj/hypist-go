package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()
	r.Use(mongoMiddleware)
	r.POST("/users", PostUser)
  r.DELETE("/users", DelUser)
	r.Run()
}

func mongoMiddleware(ctx *gin.Context) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())
	ctx.Set("db", client.Database("hypist"))
	ctx.Next()
}
