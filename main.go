package main

import (
	"context"
	"hypist/api"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.Use(mongoMiddleware)
	r.Use(cors.Default())
	r.POST("/signup", api.SignUp)
	r.DELETE("/users", verifyJWT(api.DelUser))
	r.GET("/users", verifyJWT(api.GetUser))
	r.HEAD("/users", api.LookupUser)
	r.POST("/signin", api.SignIn)
	err = r.Run()
	if err != nil {
		panic(err)
	}
}

func mongoMiddleware(ctx *gin.Context) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("DB_URI")).SetServerAPIOptions(serverAPI)
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
			Email string `json:"email"`
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, "missing authorization header")
			ctx.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if !strings.HasPrefix(authHeader, "Bearer ") || len(authParts) != 2 {
			ctx.JSON(http.StatusUnauthorized, "invalid authorization header")
			ctx.Abort()
			return
		}
		tokenString := strings.Split(authHeader, " ")[1]

		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusUnauthorized, "missing token or email")
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{"email": body.Email}, func(t *jwt.Token) (interface{}, error) {
			// TODO: import a secure key from env
			return []byte("verysecret"), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, "failed to parse token")
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, "invalid token")
			ctx.Abort()
			return
		}

		ctx.Set("email", body.Email)
		endpointHandler(ctx)
	})
}
