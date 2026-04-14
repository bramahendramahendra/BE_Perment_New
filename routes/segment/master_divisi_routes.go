package segment

import (
	handler "permen_api/domain/master_divisi/handler"
	repo "permen_api/domain/master_divisi/repo"
	service "permen_api/domain/master_divisi/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterDivisiRoutes mendaftarkan semua endpoint untuk domain Master Divisi.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-divisi/get-all → GetAllMasterDivisi   (application/json)
func MasterDivisiRoutes(r *gin.RouterGroup) {
	masterDivisiRepo := repo.NewMasterDivisiRepo(db.DB)
	masterDivisiService := service.NewMasterDivisiService(masterDivisiRepo)
	masterDivisiHandler := handler.NewMasterDivisiHandler(masterDivisiService)

	masterDivisiGroup := r.Group("master-divisi")
	masterDivisiGroup.POST("/get-all", masterDivisiHandler.GetAllMasterDivisi)
}
