package middlewares

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)
const secretKey = "your-secret-key"

// Custom middleware to check for a valid JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the request header
		tokenString := c.GetHeader("Authorization")

		// Check if the token is missing or doesn't start with "Bearer "
		if tokenString == "" || len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extract the token part (without "Bearer ")
		tokenString = tokenString[7:]

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(secretKey), nil
		})

		// Check for token parsing errors
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// If the token is valid, store the user information in the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userEmail", claims["userId"])
			c.Set("userName", claims["email"])
		}

		c.Next()
	}
}