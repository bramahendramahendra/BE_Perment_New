package segment

import (
	handler "permen_api/domain/realisasi_kpi/handler"
	repo "permen_api/domain/realisasi_kpi/repo"
	service "permen_api/domain/realisasi_kpi/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// RealisasiKpiRoutes mendaftarkan semua endpoint untuk domain realisasi KPI.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /realisasi-kpi/validate                    → ValidateRealisasiKpi                (multipart/form-data + file Excel)
//	POST /realisasi-kpi/revision                    → RevisionRealisasiKpi                (multipart/form-data + file Excel)
//	POST /realisasi-kpi/create                      → CreateRealisasiKpi                  (application/json)
//	POST /realisasi-kpi/approval                    → ApprovalRealisasiKpi                (application/json)
//	POST /realisasi-kpi/get-all-approval            → GetAllApprovalRealisasiKpi          (application/json)
//	POST /realisasi-kpi/get-all-tolakan             → GetAllTolakanRealisasiKpi           (application/json)
//	POST /realisasi-kpi/get-all-daftar-realisasi    → GetAllDaftarRealisasiKpi            (application/json)
//	POST /realisasi-kpi/get-all-daftar-approval     → GetAllDaftarApprovalRealisasiKpi    (application/json)
//	POST /realisasi-kpi/get-detail                  → GetDetailRealisasiKpi               (application/json)
func RealisasiKpiRoutes(r *gin.RouterGroup) {
	realisasiKpiRepo := repo.NewRealisasiKpiRepo(db.DB)
	realisasiKpiService := service.NewRealisasiKpiService(realisasiKpiRepo)
	realisasiKpiHandler := handler.NewRealisasiKpiHandler(realisasiKpiService)

	realisasiKpiGroup := r.Group("realisasi-kpi")
	realisasiKpiGroup.POST("/validate", realisasiKpiHandler.ValidateRealisasiKpi)
	realisasiKpiGroup.POST("/revision", realisasiKpiHandler.RevisionRealisasiKpi)
	realisasiKpiGroup.POST("/create", realisasiKpiHandler.CreateRealisasiKpi)
	realisasiKpiGroup.POST("/approval", realisasiKpiHandler.ApprovalRealisasiKpi)
	realisasiKpiGroup.POST("/get-all-approval", realisasiKpiHandler.GetAllApprovalRealisasiKpi)
	realisasiKpiGroup.POST("/get-all-tolakan", realisasiKpiHandler.GetAllTolakanRealisasiKpi)
	realisasiKpiGroup.POST("/get-all-daftar-realisasi", realisasiKpiHandler.GetAllDaftarRealisasiKpi)
	realisasiKpiGroup.POST("/get-all-daftar-approval", realisasiKpiHandler.GetAllDaftarApprovalRealisasiKpi)
	realisasiKpiGroup.POST("/get-detail", realisasiKpiHandler.GetDetailRealisasiKpi)
}
