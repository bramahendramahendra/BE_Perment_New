package service

import (
	dto "permen_api/domain/master_link_format/dto"
	repo "permen_api/domain/master_link_format/repo"
)

type (
	MasterLinkFormatServiceInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterLinkFormat digunakan oleh endpoint POST /master-perspektif/get-all.
		GetAllMasterLinkFormat() (data []dto.MasterLinkFormatResponse, err error)
	}

	masterLinkFormatService struct {
		repo repo.MasterLinkFormatRepoInterface
	}
)

func NewMasterLinkFormatService(repo repo.MasterLinkFormatRepoInterface) *masterLinkFormatService {
	return &masterLinkFormatService{repo: repo}
}
