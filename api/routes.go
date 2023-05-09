package api

import (
	"fmt"
	"hypist/database"
	"hypist/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type accessToken struct {
	Token string `json:"accessToken"`
}

func PostUser(ctx *gin.Context) {
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

func DelUser(ctx *gin.Context) {
	var request struct {
		Name string
	}

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] err: failed to delete user:\n\t%v\n", err)
		ctx.JSON(http.StatusBadRequest, fmt.Errorf("failed to delete user: %w", err))
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	err := database.DeleteUser(ctx, db, request.Name)
	if err != nil {
		fmt.Printf("[hypist] err: failed to delete user:\n\t%v\n", err)
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

	ctx.IndentedJSON(http.StatusOK, user)
}

func SignIn(ctx *gin.Context) {
	var request struct {
		Email    string
		Password string
	}

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] info: received incomplete request body:\n\t%v\n", err)
		ctx.JSON(http.StatusBadRequest, "email or password missing")
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	user, err := database.GetUser(ctx, db, "email", request.Email)
	if err != nil {
		fmt.Printf("[hypist] info: user not found:\n\t%v\n", err)
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	if err != nil {
		fmt.Printf("[hypist] warn: failed to hash password:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, "failed to hash password")
		return
	}

	correctPassword := validation.CheckPasswordHash(request.Password, user.Password)
	if !correctPassword {
		fmt.Printf("[hypist] info: incorrect password:\n\t%v\n", err)
		ctx.JSON(http.StatusUnauthorized, "incorrect password")
		return
	}

	// TODO: import a secure key from env
	secret := []byte("verysecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": request.Email,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Printf("[hypist] info: failed to sign token:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, "failed to sign token")
		return
	}

	res := accessToken{Token: tokenString}
	fmt.Println(res)
	ctx.IndentedJSON(http.StatusOK, res)
}
