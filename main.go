package main

import (
	"context"
	"fmt"
	"hypist/api"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()
	r.Use(mongoMiddleware)
	r.Use(cors.Default())
	r.POST("/signup", api.SignUp)
	r.DELETE("/users", verifyJWT(api.DelUser))
	r.GET("/users", verifyJWT(api.GetUser))
	r.HEAD("/users", api.LookupUser)
	r.POST("/signin", api.SignIn)
	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func mongoMiddleware(ctx *gin.Context) {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())
	ctx.Set("database", client.Database("hypist"))
	ctx.Next()
}

func verifyJWT(endpointHandler gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		var body struct {
      Email string  `json:"email"`
      Token string  `json:"token"`
		}

		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusUnauthorized, "missing token")
      ctx.Abort()
		}

		token, err := jwt.ParseWithClaims(body.Token, &jwt.MapClaims{"email": body.Email}, func(t *jwt.Token) (interface{}, error) {
			// TODO: import a secure key from env
			return []byte("verysecret"), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, "failed to parse token")
      ctx.Abort()
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, "invalid token")
      ctx.Abort()
		}

    ctx.Set("email", body.Email)
    endpointHandler(ctx) 
	})
}
