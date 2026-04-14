package segment

import (
	handler "permen_api/domain/master_context/handler"
	repo "permen_api/domain/master_context/repo"
	service "permen_api/domain/master_context/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterContextRoutes mendaftarkan semua endpoint untuk domain Master Context.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-context/get-all → GetAllMasterContext  (application/json)
func MasterContextRoutes(r *gin.RouterGroup) {
	masterContextRepo := repo.NewMasterContextRepo(db.DB)
	masterContextService := service.NewMasterContextService(masterContextRepo)
	masterContextHandler := handler.NewMasterContextHandler(masterContextService)

	masterContextGroup := r.Group("master-context")
	masterContextGroup.POST("/get-all", masterContextHandler.GetAllMasterContext)
}
