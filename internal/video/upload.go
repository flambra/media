package video

import (
	"log"

	"github.com/flambra/account/internal/domain"
	"github.com/flambra/account/internal/mux"
	"github.com/flambra/account/internal/s3"
	"github.com/flambra/helpers/hDb"
	"github.com/flambra/helpers/hRepository"
	"github.com/flambra/helpers/hResp"
	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	var video domain.Video
	var request domain.VideoUploadRequest

	if err := c.BodyParser(&request); err != nil {
		return hResp.BadRequestResponse(c, err.Error())
	}

	// TODO: Verificar se usuário já não upou video com mesmo nome na mesma categoria

	file, err := c.FormFile("video")
	if err != nil {
		log.Printf("Failed to get video file: %v", err)
		return hResp.BadRequestResponse(c, "Failed to get video file")
	}

	s3Url, err := s3.Upload("videos", request.UserID, request.Category, file)
	if err != nil {
		log.Printf("Failed to upload video: %v", err)
		return hResp.InternalServerErrorResponse(c, "Failed to upload video")
	}

	url, err := mux.CreateAsset(s3Url)
	if err != nil {
		log.Printf("Failed to create asset: %v", err)
		return hResp.InternalServerErrorResponse(c, "Failed to create asset")
	}

	video = domain.Video{
		UserID:      request.UserID,
		Category:    request.Category,
		Title:       request.Title,
		Description: request.Description,
		Filename:    file.Filename,
		Size:        file.Size,
		S3Url:       s3Url,
		URL:         url,
	}

	repo := hRepository.New(hDb.Get(), &video, c)
	err = repo.Create()
	if err != nil {
		return hResp.InternalServerErrorResponse(c, "Failed to save video info")
	}

	response := domain.VideoUploadResponse{
		Title:       request.Title,
		Description: request.Description,
		URL:         url,
	}

	return hResp.SuccessResponse(c, &response)
}
