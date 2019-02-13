package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
			// 账户注销
			apiV1.GET("/auth/signout", authController.Signout)

		}

	}

	return r
}
