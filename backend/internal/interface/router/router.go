package router

import (
	"moneyget/internal/domain/service"
	"moneyget/internal/interface/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	userHandler *handler.UserHandler,
	investmentHandler *handler.InvestmentHandler,
	portfolioHandler *handler.PortfolioHandler,
	jwtService service.JWTService,
) *gin.Engine {
	// Ginの本番モード設定
	gin.SetMode(gin.ReleaseMode)

	// gin.Defaultの代わりにgin.Newを使用し、必要なミドルウェアを明示的に追加
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORSミドルウェアの設定
	r.Use(corsMiddleware())

	// ルーティングの設定
	api := r.Group("/api")
	{
		// パブリックルート
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// 認証が必要なルート
		protected := api.Group("")
		protected.Use(handler.AuthMiddleware(jwtService))
		{
			// ユーザー関連
			protected.GET("/users/:id", userHandler.GetUser)

			// ポートフォリオ関連
			protected.GET("/portfolio", portfolioHandler.GetPortfolio)

			// 投資関連
			protected.POST("/investments", investmentHandler.CreateInvestment)
		}
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
