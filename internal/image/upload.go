package image

import (
	"log"
	"strconv"

	"github.com/flambra/account/internal/domain"
	"github.com/flambra/account/internal/s3"
	"github.com/flambra/helpers/hDb"
	"github.com/flambra/helpers/hRepository"
	"github.com/flambra/helpers/hResp"
	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	var image domain.Image
	var request domain.ImageUploadRequest

	if err := c.BodyParser(&request); err != nil {
		return hResp.BadRequestResponse(c, err.Error())
	}

	// TODO: Verificar se usuário já não upou video com mesmo nome na mesma categoria

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Failed to get image file: %v", err)
		return hResp.BadRequestResponse(c, "Failed to get image file")
	}

	userID, _ := strconv.ParseUint(request.UserID, 10, 64)
	url, err := s3.Upload("images", uint(userID), request.Category, file)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		return hResp.InternalServerErrorResponse(c, "Failed to upload image")
	}

	image = domain.Image{
		UserID:   uint(userID),
		Category: request.Category,
		Title:    request.Title,
		Filename: file.Filename,
		Size:     file.Size,
		URL:      url,
	}

	repo := hRepository.New(hDb.Get(), &image, c)
	err = repo.Create()
	if err != nil {
		return hResp.InternalServerErrorResponse(c, "Failed to save image info")
	}

	response := domain.ImageUploadResponse{
		Title: request.Title,
		URL:   url,
	}

	return hResp.SuccessResponse(c, &response)
}
