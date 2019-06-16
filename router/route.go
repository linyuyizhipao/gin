package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"test/controller/v1"
	"test/extend/conf"
	"test/middleware"
	"time"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(conf.ServerConf.RunMode)
	// 跨域资源共享 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  conf.CORSConf.AllowAllOrigins,
		AllowMethods:     conf.CORSConf.AllowMethods,
		AllowHeaders:     conf.CORSConf.AllowHeaders,
		ExposeHeaders:    conf.CORSConf.ExposeHeaders,
		AllowCredentials: conf.CORSConf.AllowCredentials,
		MaxAge:           conf.CORSConf.MaxAge * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//v1路由组
	apiV1 := r.Group("api/v1")
	authController := new(v1.AuthController)
	{
		// 账号注册
		apiV1.POST("/auth/signup", authController.Signup)
		// 账号登录
		apiV1.POST("/auth/signin", authController.Signin)


		//中间件作用路由组
		apiV1.Use(middleware.JWTAuth())
		{


		}
	}
	//用户抬价
	AuctionController := new(v1.AuctionController)
	{
		apiV1.POST("/auth/raise_price", AuctionController.RaisePrice)
	}



	return r
}
