package repo

import (
	dto "permen_api/domain/master_challenge/dto"
	model "permen_api/domain/master_challenge/model"

	"gorm.io/gorm"
)

type (
	MasterChallengeRepoInterface interface {
		GetAllMasterChallenge(req *dto.GetAllMasterChallengeRequest) ([]*model.MstChallenge, error)
		CheckTriwulanExists(idTriwulan string) (bool, error)
		GetDB() *gorm.DB
	}

	masterChallengeRepo struct {
		db *gorm.DB
	}
)

func NewMasterChallengeRepo(db *gorm.DB) *masterChallengeRepo {
	return &masterChallengeRepo{db: db}
}

func (r *masterChallengeRepo) GetDB() *gorm.DB {
	return r.db
}
