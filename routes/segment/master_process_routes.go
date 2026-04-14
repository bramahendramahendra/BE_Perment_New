package segment

import (
	handler "permen_api/domain/master_process/handler"
	repo "permen_api/domain/master_process/repo"
	service "permen_api/domain/master_process/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterProcessRoutes mendaftarkan semua endpoint untuk domain Master Process.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-process/get-all → GetAllMasterProcess
func MasterProcessRoutes(r *gin.RouterGroup) {
	masterProcessRepo := repo.NewMasterProcessRepo(db.DB)
	masterProcessService := service.NewMasterProcessService(masterProcessRepo)
	masterProcessHandler := handler.NewMasterProcessHandler(masterProcessService)

	masterProcessGroup := r.Group("master-process")
	masterProcessGroup.POST("/get-all", masterProcessHandler.GetAllMasterProcess)
}
