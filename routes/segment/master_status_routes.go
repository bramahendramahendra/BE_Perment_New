package segment

import (
	handler "permen_api/domain/master_status/handler"
	repo "permen_api/domain/master_status/repo"
	service "permen_api/domain/master_status/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterStatusRoutes mendaftarkan semua endpoint untuk domain Master Status.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-status/get-all → GetAllMasterStatus
//	POST /master-status/get-draft → GetDraftMasterStatus
func MasterStatusRoutes(r *gin.RouterGroup) {
	masterStatusRepo := repo.NewMasterStatusRepo(db.DB)
	masterStatusService := service.NewMasterStatusService(masterStatusRepo)
	masterStatusHandler := handler.NewMasterStatusHandler(masterStatusService)

	masterStatusGroup := r.Group("master-status")
	masterStatusGroup.POST("/get-all", masterStatusHandler.GetAllMasterStatus)
	masterStatusGroup.POST("/get-draft", masterStatusHandler.GetDraftMasterStatus)
}
