package segment

import (
	handler "permen_api/domain/edm/handler"
	service "permen_api/domain/edm/service"
	db "permen_api/pkg/database"
	"permen_api/pkg/redis"

	"github.com/gin-gonic/gin"
)

// EdmRoutes mendaftarkan semua endpoint untuk domain Edm.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /edm/kpi → GetKpi  (application/json)
func EdmRoutes(r *gin.RouterGroup) {
	edmService := service.NewEdmService(db.DB, redis.Client)
	edmHandler := handler.NewEdmHandler(edmService)

	edmGroup := r.Group("edm")
	edmGroup.POST("/kpi", edmHandler.GetKpi)
}
