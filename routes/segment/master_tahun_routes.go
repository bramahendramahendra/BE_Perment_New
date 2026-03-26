package segment

import (
	handler "permen_api/domain/master_tahun/handler"
	repo "permen_api/domain/master_tahun/repo"
	service "permen_api/domain/master_tahun/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterTahunRoutes mendaftarkan semua endpoint untuk domain Master Tahun.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-tahun/get-all → GetAllMasterTahun
func MasterTahunRoutes(r *gin.RouterGroup) {
	masterTahunRepo := repo.NewMasterTahunRepo(db.DB)
	masterTahunService := service.NewMasterTahunService(masterTahunRepo)
	masterTahunHandler := handler.NewMasterTahunHandler(masterTahunService)

	masterTahunGroup := r.Group("master-tahun")
	masterTahunGroup.POST("/get-all", masterTahunHandler.GetAllMasterTahun)
}
