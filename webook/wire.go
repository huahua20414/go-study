//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-study/webook/internal/repository"
	"go-study/webook/internal/repository/cache"
	"go-study/webook/internal/repository/dao"
	"go-study/webook/internal/service"
	"go-study/webook/internal/web"
	"go-study/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//第三方依赖
		ioc.InitDB, ioc.InitRedis, ioc.InitTencentSms, ioc.InitEmailSms,
		//初始化Dao
		dao.NewUserDao,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,

		service.NewUserService,

		web.NewUserHandler,

		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
