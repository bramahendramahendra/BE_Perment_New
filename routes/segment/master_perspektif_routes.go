package segment

import (
	handler "permen_api/domain/master_perspektif/handler"
	repo "permen_api/domain/master_perspektif/repo"
	service "permen_api/domain/master_perspektif/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterPerspektifRoutes mendaftarkan semua endpoint untuk domain Master Perspektif.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-perspektif/get-all → GetAllMasterPerspektif
func MasterPerspektifRoutes(r *gin.RouterGroup) {
	masterPerspektifRepo := repo.NewMasterPerspektifRepo(db.DB)
	masterPerspektifService := service.NewMasterPerspektifService(masterPerspektifRepo)
	masterPerspektifHandler := handler.NewMasterPerspektifHandler(masterPerspektifService)

	masterPerspektifGroup := r.Group("master-perspektif")
	masterPerspektifGroup.POST("/get-all", masterPerspektifHandler.GetAllMasterPerspektif)
}
