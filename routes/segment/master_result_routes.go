package segment

import (
	handler "permen_api/domain/master_result/handler"
	repo "permen_api/domain/master_result/repo"
	service "permen_api/domain/master_result/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterResultRoutes mendaftarkan semua endpoint untuk domain Master Result.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-result/get-all → GetAllMasterResult   (application/json)
func MasterResultRoutes(r *gin.RouterGroup) {
	masterResultRepo := repo.NewMasterResultRepo(db.DB)
	masterResultService := service.NewMasterResultService(masterResultRepo)
	masterResultHandler := handler.NewMasterResultHandler(masterResultService)

	masterResultGroup := r.Group("master-result")
	masterResultGroup.POST("/get-all", masterResultHandler.GetAllMasterResult)
}
