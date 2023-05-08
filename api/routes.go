package api

import (
	"fmt"
	"hypist/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type newUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type deleteUserRequest struct {
	Name string `json:"name"`
}

func PostUser(ctx *gin.Context) {
	var request newUserRequest

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("failed to insert user: %w", err))
		return
	}

	var runs []database.Run
	db := ctx.MustGet("database").(*mongo.Database)
	user, err := database.InsertUser(ctx, db, &database.User{Name: request.Name, Email: request.Email, Password: request.Password, Runs: runs})
	if err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to insert user: %w", err))
		return
	}

	ctx.IndentedJSON(http.StatusCreated, user)
}

func DelUser(ctx *gin.Context) {
	var request deleteUserRequest

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	err := database.DeleteUser(ctx, db, request.Name)
	if err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete user: %w", err))
	}

	ctx.IndentedJSON(http.StatusOK, "user deleted")
}

func FindUser(ctx *gin.Context) {
	email := ctx.Query("email")
	name := ctx.Query("name")

	var field string
	var value string
	if email == "" && name == "" {
		ctx.JSON(http.StatusBadRequest, "parameters name and email missing")
		return
	}

	if email == "" {
		field = "name"
		value = name
	} else {
		field = "email"
		value = email
	}

	db := ctx.MustGet("database").(*mongo.Database)
	user, err := database.GetUser(ctx, db, field, value)
	fmt.Println(user)
	if err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
}
