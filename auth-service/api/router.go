package api

import (
	"auth-service/internal/infrastructure/middleware"
	"auth-service/internal/interface/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes はAPIルートを設定
func SetupRoutes(
	router *gin.Engine,
	authController *controller.AuthController,
	authMiddleware *middleware.AuthMiddleware,
) {
	// セキュリティミドルウェアの適用
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SetCSRFToken())

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1グループ
	v1 := router.Group("/api/v1")
	{
		// 認証ルート
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/logout", authController.Logout)

			// 認証が必要なエンドポイント
			authenticated := auth.Group("")
			authenticated.Use(authMiddleware.AuthRequired())
			{
				authenticated.GET("/me", authController.Me)
			}
		}

		// 管理者専用エンドポイント
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.AuthRequired())
		admin.Use(authMiddleware.RoleRequired("admin"))
		{
			// 管理者用エンドポイントは今後追加
		}
	}
}
