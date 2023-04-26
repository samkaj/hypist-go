package main

import (
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

type deleteUserRequest struct {
	Name string `json:"name"`
}

func PostUser(ctx *gin.Context) {
	var reqBody newUserRequest

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("failed to insert user: %w", err))
		return
	}

	runs := []Run{}
	db := ctx.MustGet("db").(*mongo.Database)
	user, err := InsertUser(ctx, db, &User{Name: reqBody.Name, Email: reqBody.Email, Password: reqBody.Password, Runs: runs})
	if err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to insert user: %w", err))
		return
	}

	ctx.IndentedJSON(http.StatusCreated, user)
}

func DelUser(ctx *gin.Context) {
	var reqBody deleteUserRequest

	if err := ctx.BindJSON(&reqBody); err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
		return
	}

	db := ctx.MustGet("db").(*mongo.Database)
	err := DeleteUser(ctx, db, reqBody.Name)
	if err != nil {
		fmt.Printf("[hypist] err: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete user: %w", err))
	}
}
