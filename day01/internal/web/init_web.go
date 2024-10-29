package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-study/day01/config"
	"go-study/day01/internal/repository"
	"go-study/day01/internal/repository/dao"
	"go-study/day01/internal/service"
	"go-study/day01/internal/web/middleware"
	"go-study/day01/pkg/ginx/middleware/ratelimit"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func RegisterRoutes() *gin.Engine {
	//初始化数据库
	db := initDB()

	//初始化UserHandler
	u := initUser(db)

	//初始化路由,配置跨域和redis缓存判断用户是否登录并刷新缓存时间
	server := initServer()

	//注册用户路由
	u.RegisterUserRoutes(server)

	return server
}

func initServer() *gin.Engine {
	//初始化路由
	server := gin.Default()

	//配置限流中间件
	//初始化redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	//跨域中间件
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, // 设置允许的来源
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//允许拿什么
		AllowHeaders:  []string{"Content-Type", "Authorization", "x-jwt-token"},
		ExposeHeaders: []string{"Content-Length"},
		//是否允许携带cookie
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 缓存 CORS 设置
	}))

	//session基本使用,创建session存储在cookie中 redis存储
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"))
	//if err != nil {
	//	panic(err)
	//}
	store := memstore.NewStore([]byte("3d1c198b9d0eb074f348227c07a088bdc66910b1bb34f7678923849e45478200"))
	server.Use(sessions.Sessions("mvsession", store))
	//检验是否有session的中间件 redis实现
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	//检验是否有session的中间件 jwt实现
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())

	return server
}

// 初始化要使用的userhandler对象
func initUser(db *gorm.DB) *UserHandler {
	//初始化对象
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := NewUserHandler(svc)
	return u
}
func initDB() *gorm.DB {
	//初始化配置信息
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	//数据库同步,user表
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
