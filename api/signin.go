package api

import (
	"apskrift/database"
	"apskrift/validation"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignIn(ctx *gin.Context) {
	var request struct {
    Email    string `json:"email"`
    Password string `json:"password"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		log.Printf("received incomplete request body:\n\t%v\n", err)
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error": "email or password missing"})
		return
	}

	db := ctx.MustGet("database").(*mongo.Database)
	user, err := database.GetUser(ctx, db, "email", request.Email)
	if err != nil {
		log.Printf("user %s not found:\n\t%v\n", request.Email, err)
		ctx.JSON(http.StatusNotFound, map[string]interface{}{"error": "user not found"})
		return
	}

	if err != nil {
		log.Printf("failed to hash password:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "failed to hash password"})
		return
	}

	correctPassword := validation.CheckPasswordHash(request.Password, user.Password)
	if !correctPassword {
		log.Printf("incorrect password:\n\t%v\n", err)
		ctx.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "incorrect password"})
		return
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	expiresAt := time.Now().Add(time.Hour * 24 * 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": request.Email,
		"exp":   expiresAt,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Printf("failed to sign token:\n\t%v\n", err)
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "failed to sign token"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, map[string]interface{}{"token": tokenString, "user": user})
}
