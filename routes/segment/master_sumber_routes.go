package segment

import (
	handler "permen_api/domain/master_sumber/handler"
	repo "permen_api/domain/master_sumber/repo"
	service "permen_api/domain/master_sumber/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterSumberRoutes mendaftarkan semua endpoint untuk domain Master Sumber.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-sumber/get-all → GetAllMasterSumber   (application/json)
func MasterSumberRoutes(r *gin.RouterGroup) {
	masterSumberRepo := repo.NewMasterSumberRepo(db.DB)
	masterSumberService := service.NewMasterSumberService(masterSumberRepo)
	masterSumberHandler := handler.NewMasterSumberHandler(masterSumberService)

	masterSumberGroup := r.Group("master-sumber")
	masterSumberGroup.POST("/get-all", masterSumberHandler.GetAllMasterSumber)
}
