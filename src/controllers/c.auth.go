package controllers

import (
	"os"
	"strings"
	"time"

	"github.com/garrychrstn/kit-go/db"
	"github.com/garrychrstn/kit-go/src/helpers"
	"github.com/garrychrstn/kit-go/src/types"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthController struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewAuthController(queries *db.Queries, pool *pgxpool.Pool) *AuthController {
	return &AuthController{queries: queries, pool: pool}
}

func (q *AuthController) WhoAmI(c *gin.Context) {

}

func (q *AuthController) Login(c *gin.Context) {
	data, err := helpers.ValidateRequest[types.IRequestLogin](c)
	if err != nil {
		return
	}
	var dbUser db.User
	if strings.Contains(data.UsernameOrEmail, "@") {
		dbUser, err = q.queries.GetUserByEmail(c.Request.Context(), data.UsernameOrEmail)
	} else {
		dbUser, err = q.queries.GetUserByUsername(c.Request.Context(), data.UsernameOrEmail)
	}

	if err != nil {
		c.JSON(404, gin.H{
			"error":   "general",
			"message": "user not found",
		})
		return
	}

	if err := helpers.PasswordCompare(dbUser.Password, data.Password); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24) // Token validity for 24 hours

	claims := jwt.MapClaims{
		"user": dbUser.Email,
		"exp":  expirationTime.Unix(),
	}
	dat := gin.H{
		"id":       dbUser.ID,
		"username": dbUser.Username,
		"email":    dbUser.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token", "detail": err.Error()})
		return
	}
	c.SetCookie(
		"jwt",
		tokenString,
		int(time.Until(expirationTime).Seconds()),
		"/",
		"",
		false, // Secure: set to true in production with HTTPS
		true,  // HttpOnly: as requested
	)

	// Return data in the response body, maintaining original response structure
	c.JSON(200, gin.H{
		"ok":      true,
		"message": "Login successful",
		"data":    dat,
	})
}
