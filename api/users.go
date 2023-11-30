package api

import (
	"apskrift/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func DelUser(ctx *gin.Context) {
	var request struct {
		Email string
	}

	if err := ctx.BindJSON(&request); err != nil {
		log.Printf("failed to delete user:\n\t%v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "incomplete request body"})
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	err := database.DeleteUser(ctx, db, request.Email)
	if err != nil {
		log.Printf("failed to delete user:\n\t%v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to delete user"})
		return
	}
	ctx.IndentedJSON(http.StatusFound, map[string]interface{}{"message": "user deleted"})
}

func LookupUser(ctx *gin.Context) {
	email := ctx.Query("email")
	name := ctx.Query("name")

	var field string
	var value string
	if email == "" && name == "" {
		log.Println("received find user request without parameters")
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "parameters name and email missing"})
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
	_, err := database.GetUser(ctx, db, field, value)
	if err != nil {
		log.Printf("user not found:\n\t%v\n", err)
		ctx.JSON(http.StatusNotFound, map[string]interface{}{"error": "user not found"})
		return
	}

	ctx.IndentedJSON(http.StatusFound, map[string]interface{}{"message": "user exists"})
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
		ctx.JSON(http.StatusNotFound, map[string]interface{}{"error": "no account matches email"})
		return
	}

	ret := database.User{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Password: "ommitted",
	}

	ctx.IndentedJSON(http.StatusFound, map[string]interface{}{"user": &ret})
}
