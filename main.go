package main

import (
	"context"
	"hypist/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()
	r.Use(mongoMiddleware)
	r.Use(cors.Default())
	r.POST("/users", api.PostUser)
	r.DELETE("/users", api.DelUser)
	r.GET("/users", api.FindUser)
	r.POST("/signin", api.SignIn)
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
