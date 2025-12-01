package helpers

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func VerifyToken(string string) (*jwt.Token, jwt.MapClaims, error) {

	claims := jwt.MapClaims{}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, nil, fmt.Errorf("JWT_SECRET environment variable not set")
	}

	token, err := jwt.ParseWithClaims(string, claims, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("token parsing failed: %w", err)
	}

	// Check if the token is valid (signature, expiration, etc.)
	if !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}

	return token, claims, nil
}

func ValidateRequest[T any](c *gin.Context) (*T, error) {
	var data T
	if err := c.ShouldBindJSON(&data); err != nil {
		validationFields := make(map[string]string)

		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				field, _ := reflect.TypeOf(data).FieldByName(e.Field())
				jsonTag := field.Tag.Get("json")
				jsonField := strings.Split(jsonTag, ",")[0]

				validationFields[jsonField] = fmt.Sprintf("%s is required", jsonField)
			}
		}

		fmt.Print(err)
		c.JSON(400, gin.H{
			"error":  "validation",
			"fields": validationFields,
		})
		return nil, err
	}
	return &data, nil
}

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func PasswordCompare(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
