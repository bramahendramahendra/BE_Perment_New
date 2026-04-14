package segment

import (
	handler "permen_api/domain/master_triwulan/handler"
	repo "permen_api/domain/master_triwulan/repo"
	service "permen_api/domain/master_triwulan/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterTriwulanRoutes mendaftarkan semua endpoint untuk domain Master Triwulan.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-triwulan/get-all → GetAllTriwulan   (application/json)
func MasterTriwulanRoutes(r *gin.RouterGroup) {
	masterTriwulanRepo := repo.NewMasterTriwulanRepo(db.DB)
	masterTriwulanService := service.NewMasterTriwulanService(masterTriwulanRepo)
	masterTriwulanHandler := handler.NewMasterTriwulanHandler(masterTriwulanService)

	masterTriwulanGroup := r.Group("master-triwulan")
	masterTriwulanGroup.POST("/get-all", masterTriwulanHandler.GetAllMasterTriwulan)
}
