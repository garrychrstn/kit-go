package helpers

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

func InitR2Client() *s3.Client {
	accessKey := os.Getenv("OBJECT_ACCESS_ID")
	secretKey := os.Getenv("OBJECT_ACCESS_KEY")

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		uploadTo := os.Getenv("OBJECT_URL")
		o.BaseEndpoint = aws.String(uploadTo)
		o.UsePathStyle = true
	})
}

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

		c.JSON(400, gin.H{
			"error":  "validation",
			"fields": validationFields,
		})
		c.Abort()
		return nil, err
	}
	return &data, nil
}
