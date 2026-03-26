package routes

import (
	auth_middleware "permen_api/middleware/auth"
	"permen_api/routes/segment"

	"github.com/gin-gonic/gin"
)

func protectedRoutes(r *gin.RouterGroup) {
	r.Use(auth_middleware.BearerAuthMiddleware())

	// =============================================
	// DOMAIN: TEMPLATE
	// =============================================
	segment.TemplateRoutes(r)

	// =============================================
	// DOMAIN: PENYUSUNAN KPI
	// =============================================
	segment.PenyusunanKpiRoutes(r)

	// =============================================
	// DOMAIN: MASTER TRIWULAN
	// =============================================
	segment.MasterTriwulanRoutes(r)

	// =============================================
	// DOMAIN: MASTER PERSPEKTIF
	// =============================================
	segment.MasterPerspektifRoutes(r)

	// =============================================
	// DOMAIN: MASTER TAHUN
	// =============================================
	segment.MasterTahunRoutes(r)

	// =============================================
	// ENDPOINT TESTING — hapus setelah verified ✅
	// =============================================
	r.POST("/health", func(c *gin.Context) {
		// c.JSON(200, gin.H{
		// 	"message":      "pong - protected route berhasil",
		// 	"header_user":  c.GetHeader("User"),  // hasil dari Prioritas 6
		// 	"header_userq": c.GetHeader("userq"), // pernr | nama dari JWT
		// 	"header_pernr": c.GetHeader("pernr"), // pernr dari JWT claims
		// })

		// Ambil semua header yang ada di request
		allHeaders := make(map[string]string)
		for key, values := range c.Request.Header {
			allHeaders[key] = values[0]
		}

		c.JSON(200, gin.H{
			"message":     "pong - protected route berhasil",
			"all_headers": allHeaders,
		})
	})
	// =============================================
}
