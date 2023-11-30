package api

import (
	"apskrift/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(ctx *gin.Context) {
	var request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		log.Printf("failed to read request: %s\n", err)

		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to read request"})
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	_, err := database.InsertUser(ctx, db, &database.User{ID: primitive.NewObjectID(), Name: request.Name, Email: request.Email, Password: request.Password})
	if err != nil {
		log.Printf("failed to create user: \n\t %v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to create user"})
		return
	}

	log.Printf("created new user: \n\t %v\n", request.Name)
	ctx.IndentedJSON(http.StatusCreated, map[string]interface{}{"message": "account created"})
}
