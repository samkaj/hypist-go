package api

import (
	"fmt"
	"hypist/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(ctx *gin.Context) {
	var request struct {
		Name     string
		Email    string
		Password string
	}

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] err: %s\n", err)

		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to read request"})
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
  _, err := database.InsertUser(ctx, db, &database.User{ID: primitive.NewObjectID(), Name: request.Name, Email: request.Email, Password: request.Password})
	if err != nil {
		fmt.Printf("[hypist] err: \n\t %v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to create user"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, map[string]interface{}{"message": "account created"})
}
