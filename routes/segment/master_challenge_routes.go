package segment

import (
	handler "permen_api/domain/master_challenge/handler"
	repo "permen_api/domain/master_challenge/repo"
	service "permen_api/domain/master_challenge/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterChallengeRoutes mendaftarkan semua endpoint untuk domain Master Challenge.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-challenge/get-all → GetAllMasterChallenge
func MasterChallengeRoutes(r *gin.RouterGroup) {
	masterChallengeRepo := repo.NewMasterChallengeRepo(db.DB)
	masterChallengeService := service.NewMasterChallengeService(masterChallengeRepo)
	masterChallengeHandler := handler.NewMasterChallengeHandler(masterChallengeService)

	masterChallengeGroup := r.Group("master-challenge")
	masterChallengeGroup.POST("/get-all", masterChallengeHandler.GetAllMasterChallenge)
}
