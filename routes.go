package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
    ctx.JSON(http.StatusBadRequest, err) 
    return
  }

	runs := []Run{}
	user := User{Name: reqBody.Name, Email: reqBody.Email, Password: reqBody.Password, Runs: runs}

	// TODO: Add to db
	ctx.IndentedJSON(http.StatusCreated, user)
}
