package segment

import (
	handler "permen_api/domain/template/handler"
	repo "permen_api/domain/template/repo"
	service "permen_api/domain/template/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// TemplateRoutes mendaftarkan semua endpoint untuk domain template.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /template/format-penyusunan-kpi    → GetFormatPenyusunanKpi   (application/json body + file download)
//	POST /template/revision-penyusunan-kpi  → GetRevisionPenyusunanKpi (application/json body + file download)
//	POST /template/format-realisasi-kpi     → GetFormatRealisasiKpi    (application/json body + file download)
//	POST /template/revision-realisasi-kpi   → GetRevisionRealisasiKpi  (application/json body + file download)
func TemplateRoutes(r *gin.RouterGroup) {
	templateRepo := repo.NewTemplateRepo(db.DB)
	templateService := service.NewTemplateService(templateRepo)
	templateHandler := handler.NewTemplateHandler(templateService)

	templateGroup := r.Group("template")
	templateGroup.POST("/format-penyusunan-kpi", templateHandler.GetFormatPenyusunanKpi)
	templateGroup.POST("/revision-penyusunan-kpi", templateHandler.GetRevisionPenyusunanKpi)
	templateGroup.POST("/format-realisasi-kpi", templateHandler.GetFormatRealisasiKpi)
	templateGroup.POST("/revision-realisasi-kpi", templateHandler.GetRevisionRealisasiKpi)
}
