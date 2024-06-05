package video

// Listar url de videos de um usuário específico 
// ex: o produto que ele criou possui várias videos
// no profile dele possui um video de apresentação


// func List(c *fiber.Ctx) error {
// 	var institutions []domain.Institution
// 	institutionRepo := hRepository.New(hDb.Get(), &institutions, c)

// 	err := institutionRepo.FindAllWhere(nil)
// 	if err != nil {
// 		return hResp.InternalServerErrorResponse(c, err.Error())
// 	}

// 	var response []domain.InstitutionListResponse
// 	for _, institution := range institutions {
// 		response = append(response, domain.InstitutionListResponse{
// 			ID:   institution.ID,
// 			Name: institution.Name,
// 			Code: institution.Code,
// 			Logo: institution.Logo,
// 		})
// 	}

// 	return hResp.SuccessResponse(c, response)
// }
