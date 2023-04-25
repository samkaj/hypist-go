package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type newUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func PostUser(ctx *gin.Context) {
	var reqBody newUserRequest

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("{error: %v}", err))
		return
	}

	runs := []Run{}
	user := User{Name: reqBody.Name, Email: reqBody.Email, Password: reqBody.Password, Runs: runs}

	db := ctx.MustGet("db").(*mongo.Database)
	collection := db.Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("{error: %v}", err))
		return
	}
	ctx.IndentedJSON(http.StatusCreated, user)
}
