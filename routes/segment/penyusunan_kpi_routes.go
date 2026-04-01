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
//	POST /penyusunan-kpi/validate      → ValidatePenyusunanKpi       (multipart/form-data + file Excel)
//	POST /penyusunan-kpi/create        → CreatePenyusunanKpi         (application/json)
//	POST /penyusunan-kpi/get-all-draft → GetAllDraftPenyusunanKpi    (application/json)
//	POST /penyusunan-kpi/get-detail    → GetDetailPenyusunanKpi      (application/json)
//	POST /penyusunan-kpi/get-csv       → GetCsvPenyusunanKpi         (application/json → file download)
//	POST /penyusunan-kpi/get-pdf       → GetPdfPenyusunanKpi         (application/json → file download)
func PenyusunanKpiRoutes(r *gin.RouterGroup) {
	penyusunanKpiRepo := repo.NewPenyusunanKpiRepo(db.DB)
	penyusunanKpiService := service.NewPenyusunanKpiService(penyusunanKpiRepo)
	penyusunanKpiHandler := handler.NewPenyusunanKpiHandler(penyusunanKpiService)

	penyusunanKpiGroup := r.Group("penyusunan-kpi")
	penyusunanKpiGroup.POST("/validate", penyusunanKpiHandler.ValidatePenyusunanKpi)
	penyusunanKpiGroup.POST("/create", penyusunanKpiHandler.CreatePenyusunanKpi)
	penyusunanKpiGroup.POST("/get-all-draft", penyusunanKpiHandler.GetAllDraftPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-detail", penyusunanKpiHandler.GetDetailPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-csv", penyusunanKpiHandler.GetCsvPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-pdf", penyusunanKpiHandler.GetPdfPenyusunanKpi)
}
