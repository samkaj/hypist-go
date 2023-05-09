package api

import (
	"fmt"
	"hypist/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("failed to insert user: %w", err))
		return
	}

	var runs []database.Run
	db := ctx.MustGet("database").(*mongo.Database)
	user, err := database.InsertUser(ctx, db, &database.User{Name: request.Name, Email: request.Email, Password: request.Password, Runs: runs})
	if err != nil {
		fmt.Printf("[hypist] err: username or email taken:\n\t %v\n", err)
		ctx.JSON(http.StatusInternalServerError, "username or email taken")
		return
	}

	// TODO: import a secure key from env
	secret := []byte("verysecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": request.Email,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Printf("[hypist]: failed to sign token:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, "failed to sign token")
		return
	}

	res := struct {
		User  database.User
		Token string `json:"accessToken"`
	}{User: *user, Token: tokenString}

	ctx.IndentedJSON(http.StatusCreated, res)
}
