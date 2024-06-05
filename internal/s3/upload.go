package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var (
	s3Client *s3.Client
	bucket   string
)

func init() {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket = os.Getenv("AWS_S3_BUCKET")
	region := os.Getenv("AWS_REGION")

	if awsAccessKeyID == "" || awsSecretAccessKey == "" || bucket == "" || region == "" {
		log.Fatalf("AWS credentials, bucket name, or region are not set in environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsAccessKeyID,
			awsSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	
	log.Println("S3 client initialized")
}

func Upload(typeData string, userID uint, category string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return "", err
	}

	uuid := uuid.NewString()
	ext := ""
	if dot := strings.LastIndex(file.Filename, "."); dot != -1 {
		ext = file.Filename[dot:]
	}
	key := fmt.Sprintf("%s/%d/%s/%s%s", typeData, userID, category, uuid, ext)

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
	url = sanitizeURL(url)
	
	return url, nil
}

func sanitizeURL(url string) string {
	return strings.ReplaceAll(url, " ", "+")
}