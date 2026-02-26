package controllers

import (
	"context"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dirental/core/db"
	"github.com/dirental/core/src/helpers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TStorageController struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func StorageController(queries *db.Queries, pool *pgxpool.Pool) *TStorageController {
	return &TStorageController{queries: queries, pool: pool}
}

func (c *TStorageController) GlobalUploads() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form"})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
			return
		}

		client := helpers.InitR2Client()
		bucket := os.Getenv("OBJECT_BUCKET")
		type result struct {
			Name  string `json:"name"`
			Error string `json:"error,omitempty"`
		}

		results := make([]result, len(files))
		var wg sync.WaitGroup

		for i, header := range files {
			wg.Add(1)
			go func(i int, header *multipart.FileHeader) {
				defer wg.Done()

				file, err := header.Open()
				if err != nil {
					results[i] = result{Name: header.Filename, Error: err.Error()}
					return
				}
				defer file.Close()

				key := helpers.GenerateUUID()

				_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
					Bucket:        aws.String(bucket),
					Key:           aws.String(key),
					Body:          file,
					ContentLength: aws.Int64(header.Size),
				})
				if err != nil {
					results[i] = result{Name: header.Filename, Error: err.Error()}
					return
				}

				results[i] = result{Name: key}
			}(i, header)
		}

		wg.Wait()

		// separate successes and failures
		uploaded := []string{}
		failed := []result{}
		for _, r := range results {
			if r.Error != "" {
				failed = append(failed, r)
			} else {
				uploaded = append(uploaded, r.Name)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"uploaded": uploaded,
			"failed":   failed,
		})
	}
}
