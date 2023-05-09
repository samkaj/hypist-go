package api

import (
	"fmt"
	"hypist/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)


func DelUser(ctx *gin.Context) {
	var request struct {
		Email string
	}

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] err: failed to delete user:\n\t%v\n", err)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	err := database.DeleteUser(ctx, db, request.Email)
	if err != nil {
		fmt.Printf("[hypist] err: failed to delete user:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("failed to delete user: %w", err))
	}

	ctx.IndentedJSON(http.StatusOK, "user deleted")
}

func LookupUser(ctx *gin.Context) {
	email := ctx.Query("email")
	name := ctx.Query("name")

	var field string
	var value string
	if email == "" && name == "" {
		fmt.Println("[hypist] info: received find user request without parameters")
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
		fmt.Printf("[hypist] info: user not found:\n\t%v\n", err)
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	ctx.IndentedJSON(http.StatusFound, "user exists")
}

func GetUser(ctx *gin.Context) {
  email := ctx.GetString("email")

  if email == "" {
    ctx.JSON(http.StatusBadRequest, "no email supplied")
    return
  }

  db := ctx.MustGet("database").(*mongo.Database)
  user, err := database.GetUser(ctx, db, "email", email)
  if err != nil {
    ctx.JSON(http.StatusNotFound, "no account matches email")
    return 
  }
  
  ret := database.User{
    Email :user.Email,
    Name: user.Name,
    Password: "ommitted",
    Runs: user.Runs,
  }

  ctx.IndentedJSON(http.StatusFound, &ret)
}
