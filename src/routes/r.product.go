package routes

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dirental/core/db"
	"github.com/dirental/core/src/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(router *gin.Engine, queries *db.Queries, pool *pgxpool.Pool) {
	api := router.Group("/v1/product")
	r2 := helpers.InitR2Client()
	api.POST("/upload", uploadHandler(r2, "dirental/dev/"))

}

func uploadHandler(client *s3.Client, bucket string) gin.HandlerFunc {
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
