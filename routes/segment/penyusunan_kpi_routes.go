package segment

import (
	handler "permen_api/domain/penyusunan_kpi/handler"
	repo "permen_api/domain/penyusunan_kpi/repo"
	service "permen_api/domain/penyusunan_kpi/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// PenyusunanKpiRoutes mendaftarkan semua endpoint untuk domain penyusunan KPI.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /penyusunan-kpi/insert  → InsertKPI
func PenyusunanKpiRoutes(r *gin.RouterGroup) {
	penyusunanKpiRepo := repo.NewPenyusunanKpiRepo(db.DB)
	penyusunanKpiService := service.NewPenyusunanKpiService(penyusunanKpiRepo)
	penyusunanKpiHandler := handler.NewPenyusunanKpiHandler(penyusunanKpiService)

	penyusunanKpiGroup := r.Group("penyusunan-kpi")
	penyusunanKpiGroup.POST("/create", penyusunanKpiHandler.CreatePenyusunanKpi)
}
