package image

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/flambra/account/internal/domain"
	"github.com/flambra/helpers/hDb"
	"github.com/flambra/helpers/hRepository"
	"github.com/flambra/helpers/hResp"
	"github.com/gofiber/fiber/v2"
)

var (
	s3Client *s3.Client
	bucket   string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	bucket = os.Getenv("S3_BUCKET")
}

func Upload(c *fiber.Ctx) error {
	var request domain.ImageUploadRequest
	var image domain.Image
	repo := hRepository.New(hDb.Get(), &image, c)

	if err := c.BodyParser(&request); err != nil {
		return hResp.BadRequestResponse(c, err.Error())
	}

	if request.Title == "" {
		return hResp.BadRequestResponse(c, "Title is required")
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Failed to get image file: %v", err)
		return hResp.BadRequestResponse(c, "Failed to get image file")
	}

	url, err := uploadToS3(request.Title, file)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		return hResp.InternalServerErrorResponse(c, "Failed to upload image")
	}

	image = domain.Image{
		Title: request.Title,
		URL:   url,
	}

	err = repo.Create()
	if err != nil {
		return hResp.InternalServerErrorResponse(c, "Failed to save image info")
	}

	return hResp.SuccessResponse(c, image)
}

func uploadToS3(filename string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)

	uploader := manager.NewUploader(s3Client)
	// key := fmt.Sprintf("images/%s_%s", uuid.New().String(), file.Filename)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, filename)
	return url, nil
}
