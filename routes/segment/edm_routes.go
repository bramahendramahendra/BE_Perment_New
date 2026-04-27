package segment

import (
	handler "permen_api/domain/edm/handler"
	service "permen_api/domain/edm/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// EdmRoutes mendaftarkan semua endpoint untuk domain EDM.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /edm/realisasi → GetRealisasi (application/json)
func EdmRoutes(r *gin.RouterGroup) {
	edmService := service.NewEdmService(db.DB)
	edmHandler := handler.NewEdmHandler(edmService)

	edmGroup := r.Group("edm")
	edmGroup.POST("/realisasi", edmHandler.GetRealisasi)
}
