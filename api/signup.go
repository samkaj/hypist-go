package api

import (
	"fmt"
	"hypist/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(ctx *gin.Context) {
	var request struct {
		Name     string
		Email    string
		Password string
		Runs     []database.Run
	}

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to read request"})
		return
	}

	var runs []database.Run
	db := ctx.MustGet("database").(*mongo.Database)
	_, err := database.InsertUser(ctx, db, &database.User{Name: request.Name, Email: request.Email, Password: request.Password, Runs: runs})
	if err != nil {
		fmt.Printf("[hypist] err: \n\t %v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to create user"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, map[string]interface{}{"message": "account created"})
}
