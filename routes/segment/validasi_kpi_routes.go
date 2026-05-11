package segment

import (
	handler "permen_api/domain/validasi_kpi/handler"
	repo "permen_api/domain/validasi_kpi/repo"
	service "permen_api/domain/validasi_kpi/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// ValidasiKpiRoutes mendaftarkan semua endpoint untuk domain validasi KPI.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /validasi-kpi/input                      → InputValidasi                   (application/json) — validate + create + revision
//	POST /validasi-kpi/approve                    → ApproveValidasi                 (application/json)
//	POST /validasi-kpi/reject                     → RejectValidasi                  (application/json)
//	POST /validasi-kpi/batal                      → ValidasiBatal                   (application/json)
//	POST /validasi-kpi/get-all-approval           → GetAllApprovalValidasi          (application/json)
//	POST /validasi-kpi/get-all-tolakan            → GetAllTolakanValidasi           (application/json)
//	POST /validasi-kpi/get-all-daftar-penyusunan  → GetAllDaftarPenyusunanValidasi  (application/json)
//	POST /validasi-kpi/get-all-daftar-approval    → GetAllDaftarApprovalValidasi    (application/json)
//	POST /validasi-kpi/get-all-validasi           → GetAllValidasi                  (application/json)
//	POST /validasi-kpi/get-detail                 → GetDetailValidasiKpi            (application/json)
//	POST /validasi-kpi/get-excel                  → GetExcelValidasiKpi             (application/json) — download Excel
//	POST /validasi-kpi/get-pdf                    → GetPdfValidasiKpi               (application/json) — download PDF
func ValidasiKpiRoutes(r *gin.RouterGroup) {
	validasiKpiRepo := repo.NewValidasiKpiRepo(db.DB)
	validasiKpiService := service.NewValidasiKpiService(validasiKpiRepo)
	validasiKpiHandler := handler.NewValidasiKpiHandler(validasiKpiService)

	validasiKpiGroup := r.Group("validasi-kpi")
	validasiKpiGroup.POST("/input", validasiKpiHandler.InputValidasiKpi)
	validasiKpiGroup.POST("/approve", validasiKpiHandler.ApproveValidasiKpi)
	validasiKpiGroup.POST("/reject", validasiKpiHandler.RejectValidasiKpi)
	validasiKpiGroup.POST("/get-all-approval", validasiKpiHandler.GetAllApprovalValidasiKpi)
	validasiKpiGroup.POST("/get-all-tolakan", validasiKpiHandler.GetAllTolakanValidasiKpi)
	validasiKpiGroup.POST("/get-all-daftar-penyusunan", validasiKpiHandler.GetAllDaftarValidasiKpi)
	validasiKpiGroup.POST("/get-all-daftar-approval", validasiKpiHandler.GetAllDaftarApprovalValidasiKpi)
	validasiKpiGroup.POST("/get-all-validasi", validasiKpiHandler.GetAllValidasiKpi)
	validasiKpiGroup.POST("/get-detail", validasiKpiHandler.GetDetailValidasiKpi)
	validasiKpiGroup.POST("/get-excel", validasiKpiHandler.GetExcelValidasiKpi)
	validasiKpiGroup.POST("/get-pdf", validasiKpiHandler.GetPdfValidasiKpi)
}
