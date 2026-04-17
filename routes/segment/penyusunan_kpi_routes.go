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
//	POST /penyusunan-kpi/validate                	→ ValidatePenyusunanKpi             (multipart/form-data + file Excel)
//	POST /penyusunan-kpi/revision                	→ RevisionPenyusunanKpi             (multipart/form-data + file Excel)
//	POST /penyusunan-kpi/create                  	→ CreatePenyusunanKpi               (application/json)
//	POST /penyusunan-kpi/approve                 	→ ApprovePenyusunanKpi				(application/json)
//	POST /penyusunan-kpi/reject                  	→ RejectPenyusunanKpi				(application/json)
//	POST /penyusunan-kpi/get-all-approval        	→ GetAllApprovalPenyusunanKpi       (application/json)
//	POST /penyusunan-kpi/get-all-tolakan         	→ GetAllTolakanPenyusunanKpi        (application/json)
//	POST /penyusunan-kpi/get-all-daftar-penyusunan 	→ GetAllDaftarPenyusunanKpi         (application/json)
//	POST /penyusunan-kpi/get-all-daftar-approval 	→ GetAllDaftarApprovalPenyusunanKpi	(application/json)
//	POST /penyusunan-kpi/get-detail              	→ GetDetailPenyusunanKpi            (application/json)
//	POST /penyusunan-kpi/get-excel               	→ GetExcelPenyusunanKpi             (application/json → file download .xlsx)
//	POST /penyusunan-kpi/get-pdf                 	→ GetPdfPenyusunanKpi               (application/json → file download)
func PenyusunanKpiRoutes(r *gin.RouterGroup) {
	penyusunanKpiRepo := repo.NewPenyusunanKpiRepo(db.DB)
	penyusunanKpiService := service.NewPenyusunanKpiService(penyusunanKpiRepo)
	penyusunanKpiHandler := handler.NewPenyusunanKpiHandler(penyusunanKpiService)

	penyusunanKpiGroup := r.Group("penyusunan-kpi")
	penyusunanKpiGroup.POST("/validate", penyusunanKpiHandler.ValidatePenyusunanKpi)
	penyusunanKpiGroup.POST("/create", penyusunanKpiHandler.CreatePenyusunanKpi)
	penyusunanKpiGroup.POST("/revision", penyusunanKpiHandler.RevisionPenyusunanKpi)
	penyusunanKpiGroup.POST("/approve", penyusunanKpiHandler.ApprovePenyusunanKpi)
	penyusunanKpiGroup.POST("/reject", penyusunanKpiHandler.RejectPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-all-approval", penyusunanKpiHandler.GetAllApprovalPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-all-tolakan", penyusunanKpiHandler.GetAllTolakanPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-all-daftar-penyusunan", penyusunanKpiHandler.GetAllDaftarPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-all-daftar-approval", penyusunanKpiHandler.GetAllDaftarApprovalPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-detail", penyusunanKpiHandler.GetDetailPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-excel", penyusunanKpiHandler.GetExcelPenyusunanKpi)
	penyusunanKpiGroup.POST("/get-pdf", penyusunanKpiHandler.GetPdfPenyusunanKpi)
}
