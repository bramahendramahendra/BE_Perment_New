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
//	POST /validasi-kpi/input                    → InputValidasi                 (application/json)
//	POST /validasi-kpi/approval                 → ApprovalValidasi              (application/json)
//	POST /validasi-kpi/approve                  → ApproveValidasi               (application/json)
//	POST /validasi-kpi/reject                   → RejectValidasi                (application/json)
//	POST /validasi-kpi/batal                    → ValidasiBatal                 (application/json)
//	POST /validasi-kpi/get-all-approval         → GetAllApprovalValidasi        (application/json)
//	POST /validasi-kpi/get-all-tolakan          → GetAllTolakanValidasi         (application/json)
//	POST /validasi-kpi/get-all-daftar-penyusunan → GetAllDaftarPenyusunanValidasi (application/json)
//	POST /validasi-kpi/get-all-daftar-approval  → GetAllDaftarApprovalValidasi  (application/json)
//	POST /validasi-kpi/get-all-validasi         → GetAllValidasi                (application/json)
func ValidasiKpiRoutes(r *gin.RouterGroup) {
	validasiKpiRepo := repo.NewValidasiKpiRepo(db.DB)
	validasiKpiService := service.NewValidasiKpiService(validasiKpiRepo)
	validasiKpiHandler := handler.NewValidasiKpiHandler(validasiKpiService)

	validasiKpiGroup := r.Group("validasi-kpi")
	validasiKpiGroup.POST("/input", validasiKpiHandler.InputValidasi)
	validasiKpiGroup.POST("/approval", validasiKpiHandler.ApprovalValidasi)
	validasiKpiGroup.POST("/approve", validasiKpiHandler.ApproveValidasi)
	validasiKpiGroup.POST("/reject", validasiKpiHandler.RejectValidasi)
	validasiKpiGroup.POST("/batal", validasiKpiHandler.ValidasiBatal)
	validasiKpiGroup.POST("/get-all-approval", validasiKpiHandler.GetAllApprovalValidasi)
	validasiKpiGroup.POST("/get-all-tolakan", validasiKpiHandler.GetAllTolakanValidasi)
	validasiKpiGroup.POST("/get-all-daftar-penyusunan", validasiKpiHandler.GetAllDaftarPenyusunanValidasi)
	validasiKpiGroup.POST("/get-all-daftar-approval", validasiKpiHandler.GetAllDaftarApprovalValidasi)
	validasiKpiGroup.POST("/get-all-validasi", validasiKpiHandler.GetAllValidasi)
}
