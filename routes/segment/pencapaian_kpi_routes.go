package segment

import (
	handler "permen_api/domain/pencapaian_kpi/handler"
	repo "permen_api/domain/pencapaian_kpi/repo"
	service "permen_api/domain/pencapaian_kpi/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// PencapaianKpiRoutes mendaftarkan semua endpoint untuk domain pencapaian KPI.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /pencapaian-kpi/get-all-pencapaian  → GetAllPencapaianKpi   (application/json)
//	POST /pencapaian-kpi/get-detail          → GetDetailPencapaianKpi (application/json)
//	POST /pencapaian-kpi/get-excel           → GetExcelPencapaianKpi  (application/json) — download Excel
//	POST /pencapaian-kpi/get-pdf             → GetPdfPencapaianKpi    (application/json) — download PDF
func PencapaianKpiRoutes(r *gin.RouterGroup) {
	pencapaianKpiRepo := repo.NewPencapaianKpiRepo(db.DB)
	pencapaianKpiService := service.NewPencapaianKpiService(pencapaianKpiRepo)
	pencapaianKpiHandler := handler.NewPencapaianKpiHandler(pencapaianKpiService)

	pencapaianKpiGroup := r.Group("pencapaian-kpi")
	pencapaianKpiGroup.POST("/get-all-pencapaian", pencapaianKpiHandler.GetAllPencapaianKpi)
	pencapaianKpiGroup.POST("/get-detail", pencapaianKpiHandler.GetDetailPencapaianKpi)
	pencapaianKpiGroup.POST("/get-excel", pencapaianKpiHandler.GetExcelPencapaianKpi)
	pencapaianKpiGroup.POST("/get-pdf", pencapaianKpiHandler.GetPdfPencapaianKpi)
}
