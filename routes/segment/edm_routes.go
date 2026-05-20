package segment

import (
	handler "permen_api/domain/edm/handler"
	service "permen_api/domain/edm/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func EdmRoutes(r *gin.RouterGroup) {
	edmService := service.NewEdmService(db.DB)
	edmHandler := handler.NewEdmHandler(edmService)

	edmGroup := r.Group("edm")
	edmGroup.POST("/kpi", edmHandler.GetKpi)
}
