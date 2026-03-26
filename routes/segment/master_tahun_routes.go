package segment

import (
	handler "permen_api/domain/master_tahun/handler"
	service "permen_api/domain/master_tahun/service"

	"github.com/gin-gonic/gin"
)

// MasterTahunRoutes mendaftarkan semua endpoint untuk domain Master Tahun.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-tahun/get-all → GetAllMasterTahun
func MasterTahunRoutes(r *gin.RouterGroup) {
	masterTahunService := service.NewMasterTahunService()
	masterTahunHandler := handler.NewMasterTahunHandler(masterTahunService)

	masterTahunGroup := r.Group("master-tahun")
	masterTahunGroup.POST("/get-all", masterTahunHandler.GetAllMasterTahun)
}
