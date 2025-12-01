package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized - no token"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set claims in context
		c.Set("claims", claims)
		c.Set("user", claims["user"])
		c.Set("of_store", claims["of_store"])
		c.Set("of_role", claims["of_role"])

		c.Next()
	}
}

type Claims struct {
	User    string
	OfRole  string
	OfStore string
	Exp     int64
}

func GetClaims(c *gin.Context) (*Claims, error) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		return nil, fmt.Errorf("no claims")
	}

	mapClaims := claimsInterface.(jwt.MapClaims)

	return &Claims{
		User:    mapClaims["user"].(string),
		OfRole:  mapClaims["of_role"].(string),
		OfStore: mapClaims["of_store"].(string),
		Exp:     int64(mapClaims["exp"].(float64)),
	}, nil
}
