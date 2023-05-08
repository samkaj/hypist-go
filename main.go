package main

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()
	r.Use(mongoMiddleware)
	r.Use(cors.Default())
	r.POST("/users", PostUser)
	r.DELETE("/users", DelUser)
	r.GET("/users", FindUser)
	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func mongoMiddleware(ctx *gin.Context) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())
	ctx.Set("database", client.Database("hypist"))
	ctx.Next()
}
