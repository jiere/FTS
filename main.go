package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"fts.local/pkg/account"
	"fts.local/pkg/transaction"

	"fts.local/utils"
	_ "fts.local/docs"
)

// @title FTS API
// @version 1.0
// @description This is FTS(Funds Transfer Service) server API document.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email dickrj@163.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name token

// @BasePath /api/v1
func main() {
	utils.DBInit()
	utils.MonitorInit()	
	db := utils.DB()
	defer db.Close()

	a := account.InitMgr(db)
	t := transaction.InitMgr(db)
	a.CreateTable()
	t.CreateTable()

	router := gin.Default()
	router.Use(cors())
	router.Use(utils.Metric())

	v1 := router.Group("/api/v1")
	{
		v1.POST("/reg", a.CreateAccount)
		v1.POST("/auth", a.Login)

		acc := router.Group("/api/v1/auth/account")
		{
			acc.Use(utils.JWTAuth())
			acc.GET(":id", a.GetAccount)
			acc.GET("", a.ListAccounts)
			acc.PUT(":id", a.UpdateAccount)
			acc.DELETE(":id", a.DeleteAccount)
		}
		tx := router.Group("api/v1/auth/transaction")
		{
			tx.Use(utils.JWTAuth())
			tx.POST("", t.CreateTransaction)
			tx.GET("", t.GetTransactions)
		}

	}
	// expose metrics API to Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// The url pointing to API definition powered by Swagger
	urlString := fmt.Sprintf("http://%s:8080/swagger/doc.json", utils.Host())
	url := ginSwagger.URL(urlString)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run(":8080")
}

// For demo use, enabled such middleware to resolve CORS issues.
// For production use, need to consider more details about security.
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
