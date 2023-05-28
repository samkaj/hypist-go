package api

import (
	"fmt"
	"hypist/database"
	"hypist/validation"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignIn(ctx *gin.Context) {
	var request struct {
		Email    string
		Password string
	}

	if err := ctx.BindJSON(&request); err != nil {
		fmt.Printf("[hypist] info: received incomplete request body:\n\t%v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "email or password missing"})
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	user, err := database.GetUser(ctx, db, "email", request.Email)
	if err != nil {
		fmt.Printf("[hypist] info: user not found:\n\t%v\n", err)
		ctx.JSON(http.StatusNotFound, map[string]interface{}{"error": "user not found"})
		return
	}

	if err != nil {
		fmt.Printf("[hypist] warn: failed to hash password:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "failed to hash password"})
		return
	}

	correctPassword := validation.CheckPasswordHash(request.Password, user.Password)
	if !correctPassword {
		fmt.Printf("[hypist] info: incorrect password:\n\t%v\n", err)
		ctx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "inccorect password"})
		return
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": request.Email,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Printf("[hypist] info: failed to sign token:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "failed to sign token"})
		return
	}

  ctx.IndentedJSON(http.StatusOK, map[string]interface{}{"token": tokenString, "user": user})
}
