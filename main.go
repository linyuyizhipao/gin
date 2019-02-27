package main

import (
	"fmt"
	"test/extend/conf"
	"test/extend/logger"
	"test/extend/redis"
	"test/models"
	"test/router"
)

func main() {
	// 基本配置初始化
	conf.Setup()
	// 日志初始化
	logger.Setup()
	// 数据库初始化
	models.Setup()
	// 缓存初始化
	redis.Setup()

	//mvc路由配置
	router := router.InitRouter()

	//启动服务
	router.Run(fmt.Sprintf(":%d", conf.ServerConf.Port))
}