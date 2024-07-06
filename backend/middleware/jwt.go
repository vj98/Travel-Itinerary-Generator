package middleware

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/pkg/store"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// check for valid user, token and session
func JWTAuthMiddleware(dbClient *mongo.Client) gin.HandlerFunc {
	log.Println("checking middleware ")
	return func(c *gin.Context) {
		cfg, err := config.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*auth.CustomClaims)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token not parsed", "status": http.StatusBadRequest})
			c.Abort()
			return
		}

		t := time.Now().Local().Unix()
		if claims.ExpiresAt < t {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token is expired", "status": http.StatusBadRequest})
			c.Abort()
			return
		}

		userStore := store.NewUserStore(dbClient, "travel", "users")
		user, err := userStore.FindUserByEmail(c.Request.Context(), claims.Email_ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": fmt.Sprintf("user not registered %v", claims.Email_ID)})
			c.Abort()
			return
		}

		if user.Status != "Active" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": fmt.Sprintf("user not active contact admin %v", claims.Email_ID)})
			c.Abort()
			return
		}

		var data map[string]interface{}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		keyValue, exists := data["session"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "session not found in request"})
			c.Abort()
			return
		}

		session, ok := keyValue.(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "session is not of type string"})
			c.Abort()
			return
		}

		if user.Session != session {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusUnauthorized, "message": "session expired"})
			c.Abort()
			return
		}

		// Token is valid; you could set user information to Gin context here
		log.Println("all passed")
		// c.Set("userID", claims.Subject)
		c.Next()
	}
}
