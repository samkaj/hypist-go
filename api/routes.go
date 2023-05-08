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

type newUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type deleteUserRequest struct {
	Name string `json:"name"`
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type accessToken struct {
	Token string `json:"accessToken"`
}

type signUpResponse struct {
	User  database.User
	Token string `json:"accessToken"`
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

	res := signUpResponse{User: *user, Token: tokenString}
	fmt.Println(res)
	ctx.IndentedJSON(http.StatusCreated, res)
}

func DelUser(ctx *gin.Context) {
	var request deleteUserRequest

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] err: failed to delete user:\n\t%v\n", err)
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

func SignIn(ctx *gin.Context) {
	var request signInRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, "email or password missing")
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	user, err := database.GetUser(ctx, db, "email", request.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}
  
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, "failed to hash password")
    return
  }

	correctPassword := validation.CheckPasswordHash(request.Password, user.Password)
	if !correctPassword {
		ctx.JSON(http.StatusUnauthorized, "incorrect password")
		return
	}

  // TODO: import a secure key from env
	secret := []byte("verysecret") 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": request.Email,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
    fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, "failed to sign token")
		return
	}

	res := accessToken{Token: tokenString}
  fmt.Println(res)
	ctx.IndentedJSON(http.StatusOK, res)
}
