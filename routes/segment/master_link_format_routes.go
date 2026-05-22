package segment

import (
	handler "permen_api/domain/master_link_format/handler"
	repo "permen_api/domain/master_link_format/repo"
	service "permen_api/domain/master_link_format/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// MasterLinkFormatRoutes mendaftarkan semua endpoint untuk domain Master LinkFormat.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /master-link_format/get-all → GetAllMasterLinkFormat   (application/json)
func MasterLinkFormatRoutes(r *gin.RouterGroup) {
	masterLinkFormatRepo := repo.NewMasterLinkFormatRepo(db.DB)
	masterLinkFormatService := service.NewMasterLinkFormatService(masterLinkFormatRepo)
	masterLinkFormatHandler := handler.NewMasterLinkFormatHandler(masterLinkFormatService)

	masterLinkFormatGroup := r.Group("master-link_format")
	masterLinkFormatGroup.POST("/get-all", masterLinkFormatHandler.GetAllMasterLinkFormat)
}
