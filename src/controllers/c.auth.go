package controllers

import (
	"log"
	"os"
	"time"

	"github.com/dirental/core/db"
	"github.com/dirental/core/src/helpers"
	helperdb "github.com/dirental/core/src/helpers/db"
	"github.com/dirental/core/src/types/requests"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IAuthController struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func AuthController(queries *db.Queries, pool *pgxpool.Pool) *IAuthController {
	return &IAuthController{queries: queries, pool: pool}
}

func (q *IAuthController) Register(c *gin.Context) {
	data, err := helpers.ValidateRequest[requests.IRequestRegister](c)

	if err != nil {
		return
	}

	if data.ConfirmPassword != data.Password {
		c.JSON(400, gin.H{
			"error": "validation",
			"fields": map[string]string{
				"confirm_password": "Passwords do not match",
			},
		})
		return
	}
	hashed, err := helperdb.PasswordHash(data.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "general",
			"message": "failed to hash password",
		})
		return
	}
	user := db.RegisterUserParams{
		Email:       data.Email,
		Password:    hashed,
		Name:        data.Name,
		Username:    data.Username,
		PhoneNumber: pgtype.Text{String: data.PhoneNumber, Valid: true},
	}

	existingUser, err := q.queries.GetUserByAny(c.Request.Context(), db.GetUserByAnyParams{
		Email: data.Email, Username: data.Username, PhoneNumber: pgtype.Text{String: data.PhoneNumber, Valid: true},
	})

	if existingUser.ID.Valid {
		c.JSON(403, gin.H{
			"error": "validation",
			"fields": map[string]string{
				"email":    "Email already exists",
				"username": "Username already exists",
				"phone":    "Phone number already exists",
			},
			"user": existingUser,
		})
		return
	}

	dbUser, err := q.queries.RegisterUser(c.Request.Context(), user)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error":   "general",
			"message": "failed to register user",
		})
		return
	}

	c.JSON(201, gin.H{
		"ok":      true,
		"message": "User registered successfully",
		"data": gin.H{
			"id":       dbUser.ID,
			"username": dbUser.Username,
			"email":    dbUser.Email,
		},
	})
}

func (q *IAuthController) Login(c *gin.Context) {
	data, err := helpers.ValidateRequest[requests.IRequestLogin](c)
	if err != nil {
		return
	}
	var dbUser db.User
	dbUser, err = q.queries.GetUserByAny(
		c.Request.Context(),
		db.GetUserByAnyParams{
			Email:       data.UsernameOrEmail,
			Username:    data.UsernameOrEmail,
			PhoneNumber: helperdb.SafeString(&data.UsernameOrEmail),
		})

	if err != nil {
		c.JSON(404, gin.H{
			"error":   "general",
			"message": "user not found",
		})
		return
	}

	if err := helperdb.PasswordCompare(dbUser.Password, data.Password); err != nil {
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

	c.JSON(200, gin.H{
		"ok":      true,
		"message": "Login successful",
		"data":    dat,
	})
}

func (q *IAuthController) Logout(c *gin.Context) {
	c.SetCookie(
		"jwt",
		"", // Empty token string to clear the cookie
		-1, // MaxAge set to -1 to delete the cookie
		"/",
		"",
		false, // Secure: set to true in production with HTTPS
		true,  // HttpOnly: as requested
	)

	c.JSON(200, gin.H{
		"ok":      true,
		"message": "Logout successful",
	})
}
