package segment

import (
	handler "permen_api/domain/template/handler"
	service "permen_api/domain/template/service"

	"github.com/gin-gonic/gin"
)

// TemplateRoutes mendaftarkan semua endpoint untuk domain template.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /template/format-penyusunan-kpi → GetFormatPenyusunanKpi (application/json body + file download response)
func TemplateRoutes(r *gin.RouterGroup) {
	templateService := service.NewTemplateService()
	templateHandler := handler.NewTemplateHandler(templateService)

	templateGroup := r.Group("template")
	templateGroup.POST("/format-penyusunan-kpi", templateHandler.GetFormatPenyusunanKpi)
}
