package segment

import (
	handler "permen_api/domain/master_method/handler"
	repo "permen_api/domain/master_method/repo"
	service "permen_api/domain/master_method/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterMethodRoutes mendaftarkan semua endpoint untuk domain Master Method.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-method/get-all → GetAllMasterMethod
func MasterMethodRoutes(r *gin.RouterGroup) {
	masterMethodRepo := repo.NewMasterMethodRepo(db.DB)
	masterMethodService := service.NewMasterMethodService(masterMethodRepo)
	masterMethodHandler := handler.NewMasterMethodHandler(masterMethodService)

	masterMethodGroup := r.Group("master-method")
	masterMethodGroup.POST("/get-all", masterMethodHandler.GetAllMasterMethod)
}
