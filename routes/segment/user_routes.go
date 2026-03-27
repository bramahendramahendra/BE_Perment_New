package segment

import (
	handler "permen_api/domain/user/handler"
	repo "permen_api/domain/user/repo"
	service "permen_api/domain/user/service"
	db "permen_api/pkg/database"

	"github.com/gin-gonic/gin"
)

// UserRoutes mendaftarkan semua endpoint untuk domain User.
// Endpoint ini berada di bawah protected route (memerlukan Bearer Auth).
//
// Daftar endpoint:
//
//	POST /user/get-all → GetAllUser
func UserRoutes(r *gin.RouterGroup) {
	userRepo := repo.NewUserRepo(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	userGroup := r.Group("user")
	userGroup.POST("/get-all", userHandler.GetAllUser)
}
