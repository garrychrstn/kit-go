package helpers

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type TStorageHandler struct {
	bucket string
	client *s3.Client
}

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

func UploadHandler(client *s3.Client, bucket string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
			return
		}
		defer file.Close()

		_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:        aws.String(bucket),
			Key:           aws.String(header.Filename),
			Body:          file,
			ContentLength: aws.Int64(header.Size),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"uploaded": header.Filename})
	}
}

