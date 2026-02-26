package helpers

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

type StorageHandler struct {
	bucket string
	client *s3.Client
}
