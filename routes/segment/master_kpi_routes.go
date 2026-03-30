package segment

import (
	handler "permen_api/domain/master_kpi/handler"
	repo "permen_api/domain/master_kpi/repo"
	service "permen_api/domain/master_kpi/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterKpiRoutes mendaftarkan semua endpoint untuk domain Master KPI.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-kpi/get-all → GetAllMasterKpi
func MasterKpiRoutes(r *gin.RouterGroup) {
	masterKpiRepo := repo.NewMasterKpiRepo(db.DB)
	masterKpiService := service.NewMasterKpiService(masterKpiRepo)
	masterKpiHandler := handler.NewMasterKpiHandler(masterKpiService)

	masterKpiGroup := r.Group("master-kpi")
	masterKpiGroup.POST("/get-all", masterKpiHandler.GetAllMasterKpi)
}
